package jupitergames

import (
	"bufio"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	colly "github.com/gocolly/colly/v2"
	queue "github.com/gocolly/colly/v2/queue"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	http "github.com/hashicorp/go-retryablehttp"

	"github.com/kodabb/go-mtgban/mtgban"
	"github.com/kodabb/go-mtgban/mtgdb"
	"github.com/kodabb/go-mtgban/mtgjson"
)

const (
	defaultConcurrency = 8

	jupInventoryURL = "https://jupitergames.info/store/start"
	jupBuylistURL   = "https://jupitergames.info/store/find/buypriceall"
)

type Jupitergames struct {
	LogCallback    mtgban.LogCallbackFunc
	InventoryDate  time.Time
	BuylistDate    time.Time
	MaxConcurrency int

	inventory mtgban.InventoryRecord
	buylist   mtgban.BuylistRecord
}

func NewScraper() *Jupitergames {
	jup := Jupitergames{}
	jup.inventory = mtgban.InventoryRecord{}
	jup.buylist = mtgban.BuylistRecord{}
	jup.MaxConcurrency = defaultConcurrency
	return &jup
}

type responseChan struct {
	card     *mtgdb.Card
	invEntry *mtgban.InventoryEntry
}

func (jup *Jupitergames) printf(format string, a ...interface{}) {
	if jup.LogCallback != nil {
		jup.LogCallback("[JUP] "+format, a...)
	}
}

func (jup *Jupitergames) scrape() error {
	channel := make(chan responseChan)

	c := colly.NewCollector(
		colly.AllowedDomains("jupitergames.info"),

		colly.CacheDir(fmt.Sprintf(".cache/%d", time.Now().YearDay())),

		colly.Async(true),
	)

	c.SetClient(cleanhttp.DefaultClient())

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
		Parallelism: jup.MaxConcurrency,
	})

	c.OnRequest(func(r *colly.Request) {
		//jup.printf("Visiting %s", r.URL.String())
	})

	c.OnHTML(`form[class="form"]`, func(e *colly.HTMLElement) {
		id := e.ChildAttr(`input[name="id"]`, "value")
		cardName := e.ChildAttr(`input[name="name"]`, "value")
		priceStr := e.ChildAttr(`input[name="price"]`, "value")
		variant := e.ChildAttr(`input[name="variant"]`, "value")
		edition := e.ChildAttr(`input[name="set"]`, "value")
		format := e.ChildAttr(`input[name="format"]`, "value")
		language := e.ChildAttr(`input[name="language"]`, "value")
		conditions := e.ChildAttr(`input[name="condition"]`, "value")
		qtyStr := e.ChildAttr(`input[name="qty"]`, "max")

		// Skip promo cards sneaking in normal sets
		if !strings.HasSuffix(e.Request.URL.Path, url.QueryEscape(edition)) {
			return
		}

		switch conditions {
		case "NM", "SP":
		default:
			jup.printf("Unsupported %s condition for %s %s", conditions, cardName, edition)
			return
		}

		if !strings.Contains(language, "English") {
			jup.printf("%s %s %s", cardName, edition)
			return
		}

		qty, err := strconv.Atoi(qtyStr)
		if err != nil {
			jup.printf("%s %s %v", cardName, edition, err)
			return
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			jup.printf("%s %s %v", cardName, edition, err)
			return
		}

		if price == 0.0 || qty == 0 {
			return
		}

		theCard, err := preprocess(cardName, variant, edition, format)
		if err != nil {
			return
		}
		cc, err := theCard.Match()
		if err != nil {
			switch theCard.Edition {
			default:
				jup.printf("%q", theCard)
				jup.printf("%v", err)
			}
			return
		}

		var out responseChan
		out = responseChan{
			card: cc,
			invEntry: &mtgban.InventoryEntry{
				Conditions: conditions,
				Price:      price,
				Quantity:   qty,
				URL:        "https://jupitergames.info/store/find/id/" + id,
			},
		}

		channel <- out
	})

	q, _ := queue.New(
		jup.MaxConcurrency,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	resp, err := http.Get(jupInventoryURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	doc.Find(`script`).Each(func(_ int, s *goquery.Selection) {
		script := s.Text()
		if !strings.Contains(script, "availableSets") {
			return
		}
		parseEds := false
		scanner := bufio.NewScanner(strings.NewReader(script))
		i := 0
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "availableSets") {
				parseEds = true
				continue
			} else if strings.Contains(line, "];") {
				break
			}
			if !parseEds {
				continue
			}
			line = strings.Replace(line, "\"", "", -1)
			line = strings.Replace(line, ",", "", 1)
			line = strings.TrimSpace(line)

			if line == "Alternate Fourth Edition" {
				continue
			}
			i++
			link := "https://jupitergames.info/store/find/set/" + url.QueryEscape(line)
			q.AddURL(link)
		}
		jup.printf("Parsing %d editions", i)
	})

	q.Run(c)

	go func() {
		c.Wait()
		close(channel)
	}()

	for res := range channel {
		err := jup.inventory.Add(res.card, res.invEntry)
		if err != nil {
			jup.printf("%v", err)
		}
	}

	jup.InventoryDate = time.Now()

	return nil
}

func (jup *Jupitergames) Inventory() (mtgban.InventoryRecord, error) {
	if len(jup.inventory) > 0 {
		return jup.inventory, nil
	}

	start := time.Now()
	jup.printf("Inventory scraping started at %s", start)

	err := jup.scrape()
	if err != nil {
		return nil, err
	}
	jup.printf("Inventory scraping took %s", time.Since(start))

	return jup.inventory, nil
}

func (jup *Jupitergames) parseBL() error {
	resp, err := http.Get(jupBuylistURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	list := doc.Find(`div[class="row-fluid"]`).Has("span[style='color:red']").Text()
	scanner := bufio.NewScanner(strings.NewReader(list))
	i := -1
	for scanner.Scan() {
		i++
		line := scanner.Text()

		// Replace the split cards separator
		if strings.Contains(line, "||") {
			line = strings.Replace(line, " || ", " // ", 1)
		}

		// Split as if it was a CSV
		record := strings.Split(line, "|")

		// For the first and last line which is empty
		if len(record) <= 1 {
			continue
		}
		// Skip the second line, which is the hader
		cardName := strings.TrimSpace(record[0])
		if cardName == "NAME" {
			continue
		}

		// Make sure that the split card don't mess up our record splitting
		simple, err := mtgdb.CardSimple(cardName)
		if err == nil && cardName != "Bind" {
			if len(simple.Names) > 1 && mtgjson.NormContains(record[1], simple.Names[1]) {
				line = strings.Replace(line, " | ", " // ", 1)
				record = strings.Split(line, "|")
				cardName = strings.TrimSpace(record[0])
			}
		}

		edition := strings.TrimSpace(record[1])

		priceStr := strings.TrimSpace(record[len(record)-1])

		// Foil and variant are optional, and not very well identifiable
		notes := strings.TrimSpace(record[len(record)-2])
		format := ""
		variant := ""
		if notes == "*FOIL*" {
			format = notes
		} else if notes != edition {
			variant = notes
		}
		if variant == "" && len(record) > 3 {
			maybe := strings.TrimSpace(record[len(record)-3])
			if maybe != edition && maybe != cardName {
				variant = maybe
			}
		}

		priceStr = strings.Replace(priceStr, "$", "", 1)
		priceStr = strings.Replace(priceStr, ",", "", 1)
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			jup.printf(err.Error())
			continue
		}

		// No difference from the buylist page
		if cardName == "Demonic Tutor" && edition == "Promo - Judge Gift Program" {
			edition = "J20"
			if price > 100 {
				edition = "G08"
			}
		} else if cardName == "Path to Exile" && edition == "Promo - MagicFest 2020" {
			if price > 10 {
				format = "*FOIL*"
			}
		}

		theCard, err := preprocess(cardName, variant, edition, format)
		if err != nil {
			continue
		}
		cc, err := theCard.Match()
		if err != nil {
			jup.printf("%q", theCard)
			jup.printf("%v", err)
			continue
		}

		var priceRatio, sellPrice float64

		invCards := jup.inventory[*cc]
		for _, invCard := range invCards {
			sellPrice = invCard.Price
			break
		}
		if sellPrice > 0 {
			priceRatio = price / sellPrice * 100
		}

		buyEntry := &mtgban.BuylistEntry{
			Quantity:   1,
			BuyPrice:   price,
			TradePrice: price * 1.25,
			PriceRatio: priceRatio,
			URL:        "https://jupitergames.info/store/find/buypricebyname/" + url.QueryEscape(cardName),
		}

		err = jup.buylist.Add(cc, buyEntry)
		if err != nil {
			jup.printf(err.Error())
			continue
		}
	}

	return nil
}

func (jup *Jupitergames) Buylist() (mtgban.BuylistRecord, error) {
	if len(jup.buylist) > 0 {
		return jup.buylist, nil
	}

	start := time.Now()
	jup.printf("Buylist scraping started at %s", start)

	err := jup.parseBL()
	if err != nil {
		return nil, err
	}
	jup.printf("Buylist scraping took %s", time.Since(start))

	return jup.buylist, nil
}

func (jup *Jupitergames) Grading(card mtgdb.Card, entry mtgban.BuylistEntry) (grade map[string]float64) {
	grade = map[string]float64{
		"SP": 0.7, "MP": 0.5, "HP": 0.3,
	}
	return
}

func (jup *Jupitergames) Info() (info mtgban.ScraperInfo) {
	info.Name = "Jupiter Games"
	info.Shorthand = "JUP"
	info.InventoryTimestamp = jup.InventoryDate
	info.BuylistTimestamp = jup.BuylistDate
	return
}
