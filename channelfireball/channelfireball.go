package channelfireball

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	colly "github.com/gocolly/colly/v2"
	cleanhttp "github.com/hashicorp/go-cleanhttp"

	"github.com/kodabb/go-mtgban/mtgban"
	"github.com/kodabb/go-mtgban/mtgjson"
)

const (
	maxConcurrency  = 8
	cfbInventoryURL = "https://store.channelfireball.com/catalog/magic_singles/8"
	cfbBuylistURL   = "https://store.channelfireball.com/buylist/magic_singles/8"

	modeInventory = "inventory"
	modeBuylist   = "buylist"
)

type Channelfireball struct {
	LogCallback mtgban.LogCallbackFunc

	db        mtgjson.MTGDB
	inventory map[string][]mtgban.InventoryEntry
	buylist   map[string]mtgban.BuylistEntry

	norm *mtgban.Normalizer
}

func NewScraper(db mtgjson.MTGDB) *Channelfireball {
	cfb := Channelfireball{}
	cfb.db = db
	cfb.inventory = map[string][]mtgban.InventoryEntry{}
	cfb.buylist = map[string]mtgban.BuylistEntry{}
	cfb.norm = mtgban.NewNormalizer()
	return &cfb
}

type resultChan struct {
	err  error
	card cfbCard
}

func (cfb *Channelfireball) printf(format string, a ...interface{}) {
	if cfb.LogCallback != nil {
		cfb.LogCallback("[CFB] "+format, a...)
	}
}

func (cfb *Channelfireball) scrape(mode string) error {
	channel := make(chan cfbCard)

	c := colly.NewCollector(
		colly.AllowedDomains("store.channelfireball.com"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted - daily
		colly.CacheDir(fmt.Sprintf(".cache/%d", time.Now().YearDay())),

		colly.Async(true),
	)

	c.SetClient(cleanhttp.DefaultClient())

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
		Parallelism: maxConcurrency,
	})

	c.OnRequest(func(r *colly.Request) {
		q := r.URL.Query()
		if q.Get("page") == "" {
			cfb.printf("Visiting %s", r.URL.String())
		}
	})

	// Callback for links on scraped pages (edition names)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Consider only "/{buylist,catalog}/magic_singles-<catergory>-<edition>/<id>" links
		linkDepth := strings.Count(link, "/")
		ok := (mode == modeInventory && strings.HasPrefix(link, "/catalog")) ||
			(mode == modeBuylist && strings.HasPrefix(link, "/buylist"))
		if linkDepth == 3 && ok {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnHTML("form[class='add-to-cart-form']", func(e *colly.HTMLElement) {
		// Skip out of stock items
		dataVid := e.Attr("data-vid")
		if dataVid == "" {
			return
		}

		priceStr := e.Attr("data-price")
		priceStr = strings.Replace(priceStr, "$", "", 1)
		priceStr = strings.Replace(priceStr, ",", "", 1)
		cardPrice, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			cfb.printf("%v", err)
			return
		}

		edition := e.Attr("data-category")
		// Strip stars indicating edition in preorder
		edition = strings.Replace(edition, "*", "", -1)

		// Skip non-playable and oversized card sets
		switch {
		case edition == "Promos: Hero's Path" ||
			(strings.Contains(edition, "{") && strings.Contains(edition, "}")):
			return
		}

		cond := e.Attr("data-variant")
		fields := strings.Split(cond, ", ")
		cond = fields[0]
		if len(fields) > 1 && fields[1] != "English" {
			return
		}
		switch cond {
		case "NM-Mint":
			cond = "NM"
		case "Slightly Played":
			cond = "SP"
		case "Moderately Played":
			cond = "MP"
		case "Damaged":
			cond = "HP"
		default:
			cfb.printf("Unsupported %s condition", cond)
			return
		}

		cardName := e.Attr("data-name")
		if cardName == "" {
			// Quotes are not escaped
			return
		}
		// Skip duplicated cards, not yet tracked
		switch cardName {
		case "Blaze (Alternate Art - Deck)", "Blaze (Alternate Art - Booster)",
			"Crystalline Sliver - Arena 2003":
			return
		}

		cardName = strings.Replace(cardName, "–", "-", -1)

		// Skip tokens and similar cards
		if strings.Contains(cardName, "Token") || strings.Contains(cardName, "token") ||
			strings.Contains(cardName, "Checklist") ||
			strings.Contains(cardName, "Filler") ||
			strings.Contains(cardName, "APAC Land Set") ||
			strings.Contains(cardName, "Emblem") {
			return
		}
		switch cardName {
		case "Experience Counter", "Poison Counter", "Experience Card",
			"Goblin", "Pegasus", "Sheep", "Soldier", "Squirrel", "Zombie",
			"Standard Placeholder", "Blank Card", "Splendid Genesis",
			"Black ": // Black "M" Filler Card
			return
		}

		// skip non-english versions of this card
		if strings.HasPrefix(cardName, "Mana Crypt (Book Promo) (") {
			return
		}

		isFoil := false
		if strings.Contains(cardName, " Foil") ||
			// Our Market Research Shows that really long names hide card properties
			dataVid == "1061099" || dataVid == "297099" {
			isFoil = true
		}

		// Drop pointeless tags
		cardName = strings.Replace(cardName, " - Foil", "", 1)
		cardName = strings.Replace(cardName, " - Hero's Path", "", 1)
		cardName = strings.Replace(cardName, " (Masterpiece Foil)", "", 1)

		// Correctly put variants in the correct tag (within parenthesis)
		tags := []string{
			"Magic League Promo", "Draft Weekend Promo", "Draft Weekend",
			"Planeswalker Weekend Promo", "Media Promo", "Open House Promo",
			"Bundle Promo", "SDCC 2019 Exclusive", "FNM 2017", "FNM 2019",
			"FNM Promo 2019", "DCI Judge Promo", "Judge Academy Promo",
			"Buy-a-Box Promo", "Store Championship Promo", "Dark Frame Promo",
			"Treasure Map",
		}
		for _, tag := range tags {
			cardName = strings.Replace(cardName, " "+tag, " ("+tag+")", 1)
		}

		// Make sure that variants are separated from the name
		parIndex := strings.Index(cardName, "(")
		if parIndex-1 > 0 && parIndex-1 < len(cardName) && cardName[parIndex-1] != ' ' {
			cardName = strings.Replace(cardName, "(", " (", 1)
		}

		// Split by () and by -, rebuild the cardname in a standardized way
		fields = mtgban.SplitVariants(cardName)
		subfields := strings.Split(fields[0], " - ")
		cardName = subfields[0]
		for _, field := range fields[1:] {
			field = strings.Replace(field, " - ", " ", -1)
			cardName += " (" + field + ")"
		}
		for _, field := range subfields[1:] {
			cardName += " (" + field + ")"
		}
		cardName = strings.Replace(cardName, " - ", " ", 1)

		// Fixup any expected errors
		lutName, found := cardTable[cardName]
		if found {
			cardName = lutName
		}

		qty := 0
		e.ForEach("input", func(_ int, elem *colly.HTMLElement) {
			if elem.Attr("class") == "qty" {
				qty, err = strconv.Atoi(elem.Attr("max"))
				if err != nil {
					return
				}
			}
		})
		if err != nil {
			cfb.printf("%v", err)
			return
		}

		card := cfbCard{
			Key:        dataVid,
			Name:       cardName,
			Edition:    edition,
			Foil:       isFoil,
			Conditions: cond,
			Price:      cardPrice,
			Quantity:   qty,
		}

		channel <- card
	})

	if mode == modeInventory {
		c.Visit(cfbInventoryURL)
	} else if mode == modeBuylist {
		c.Visit(cfbBuylistURL)
	} else {
		return fmt.Errorf("Unsupported mode %s", mode)
	}

	go func() {
		c.Wait()
		close(channel)
	}()

	// The same pattern is repeated exactly 3 times, store the simple key for
	// the processed cards and skip the duplicatoin
	processed := map[string]bool{}

	for card := range channel {
		if processed[card.Key] {
			continue
		}
		processed[card.Key] = true

		cc, err := cfb.convert(&card)
		if err != nil {
			switch {
			// Ignore errors coming from lands from these two editions only
			case strings.HasPrefix(card.Name, "Plains") ||
				strings.HasPrefix(card.Name, "Island") ||
				strings.HasPrefix(card.Name, "Swamp") ||
				strings.HasPrefix(card.Name, "Mountain") ||
				strings.HasPrefix(card.Name, "Forest"):
				if card.Edition != "5th Edition" && card.Edition != "Gift Boxes: Battle Royale" {
					cfb.printf("%v", err)
				}
			default:
				cfb.printf("%v", err)
			}
			continue
		}

		if mode == modeInventory {
			if card.Quantity > 0 && card.Price > 0 {
				out := mtgban.InventoryEntry{
					Card:       *cc,
					Conditions: card.Conditions,
					Price:      card.Price,
					Quantity:   card.Quantity,
				}
				err := cfb.InventoryAdd(out)
				if err != nil {
					switch cc.Name {
					// Ignore errors coming from lands for now
					case "Plains", "Island", "Swamp", "Mountain", "Forest":
					default:
						cfb.printf("%v", err)
					}
					continue
				}
			}
		}
		if mode == modeBuylist {
			if card.Quantity > 0 && card.Price > 0 && card.Conditions == "NM" {
				var sellPrice, priceRatio, qtyRatio float64
				sellQty := 0

				invCards := cfb.inventory[cc.Id]
				for _, invCard := range invCards {
					if invCard.Conditions == card.Conditions {
						sellPrice = invCard.Price
						sellQty = invCard.Quantity
						break
					}
				}

				if sellPrice > 0 {
					priceRatio = card.Price / sellPrice * 100
				}
				if sellQty > 0 {
					qtyRatio = float64(card.Quantity) / float64(sellQty) * 100
				}

				out := mtgban.BuylistEntry{
					Card:          *cc,
					Conditions:    card.Conditions,
					BuyPrice:      card.Price,
					TradePrice:    card.Price * 1.3,
					Quantity:      card.Quantity,
					PriceRatio:    priceRatio,
					QuantityRatio: qtyRatio,
				}
				err := cfb.BuylistAdd(out)
				if err != nil {
					switch cc.Name {
					// Ignore errors coming from lands for now
					case "Plains", "Island", "Swamp", "Mountain", "Forest":
					default:
						cfb.printf("%v", err)
					}
					continue
				}
			}
		}
	}

	return nil
}

func (cfb *Channelfireball) InventoryAdd(card mtgban.InventoryEntry) error {
	entries, found := cfb.inventory[card.Id]
	if found {
		for _, entry := range entries {
			if entry.Conditions == card.Conditions && entry.Price == card.Price {
				return fmt.Errorf("Attempted to add a duplicate inventory card:\n-new: %v\n-old: %v", card, entry)
			}
		}
	}

	cfb.inventory[card.Id] = append(cfb.inventory[card.Id], card)
	return nil
}

func (cfb *Channelfireball) Inventory() (map[string][]mtgban.InventoryEntry, error) {
	if len(cfb.inventory) > 0 {
		return cfb.inventory, nil
	}

	cfb.printf("Empty inventory, scraping started")

	err := cfb.scrape(modeInventory)
	if err != nil {
		return nil, err
	}

	return cfb.inventory, nil
}

func (cfb *Channelfireball) BuylistAdd(card mtgban.BuylistEntry) error {
	entry, found := cfb.buylist[card.Id]
	if found {
		if entry.Conditions == card.Conditions && entry.BuyPrice == card.BuyPrice {
			return fmt.Errorf("Attempted to add a duplicate buylist card:\n-new: %v\n-old: %v", card, entry)
		}
	}

	cfb.buylist[card.Id] = card
	return nil
}

func (cfb *Channelfireball) Buylist() (map[string]mtgban.BuylistEntry, error) {
	if len(cfb.buylist) > 0 {
		return cfb.buylist, nil
	}

	cfb.printf("Empty buylist, scraping started")

	err := cfb.scrape(modeBuylist)
	if err != nil {
		return nil, err
	}

	return cfb.buylist, nil
}

var fourHorsemenDate = time.Date(1993, time.August, 1, 0, 0, 0, 0, time.UTC)
var premodernDate = time.Date(1994, time.August, 1, 0, 0, 0, 0, time.UTC)
var modernDate = time.Date(2003, time.July, 1, 0, 0, 0, 0, time.UTC)

func (cfb *Channelfireball) Grading(entry mtgban.BuylistEntry) (grade map[string]float64) {
	var setDate time.Time
	for _, set := range cfb.db {
		if set.Name == entry.Card.Set {
			setDate, _ = time.Parse("2006-01-02", set.ReleaseDate)
			break
		}
	}

	switch {
	case entry.Card.Foil:
		grade = map[string]float64{
			"SP": 0.7, "MP": 0.5, "HP": 0.3,
		}
	case setDate.After(fourHorsemenDate) && setDate.Before(premodernDate.AddDate(0, 0, -1)):
		grade = map[string]float64{
			"SP": 0.5, "MP": 0.25, "HP": 0.1,
		}
	case setDate.After(premodernDate) && setDate.Before(modernDate.AddDate(0, 0, -1)):
		grade = map[string]float64{
			"SP": 0.7, "MP": 0.5, "HP": 0.3,
		}
	case setDate.After(modernDate):
		grade = map[string]float64{
			"SP": 0.8, "MP": 0.6, "HP": 0.4,
		}
	}

	return
}

func (cfb *Channelfireball) Info() (info mtgban.ScraperInfo) {
	info.Name = "Channel Fireball"
	info.Shorthand = "CFB"
	return
}
