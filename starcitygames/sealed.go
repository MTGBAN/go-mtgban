package starcitygames

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mtgban/go-mtgban/mtgban"
	"github.com/mtgban/go-mtgban/mtgmatcher"
)

type StarcitygamesSealed struct {
	LogCallback    mtgban.LogCallbackFunc
	inventoryDate  time.Time
	buylistDate    time.Time
	MaxConcurrency int

	Affiliate string

	inventory mtgban.InventoryRecord
	buylist   mtgban.BuylistRecord

	client *SCGClient
}

func NewScraperSealed(guid, bearer string) *StarcitygamesSealed {
	scg := StarcitygamesSealed{}
	scg.inventory = mtgban.InventoryRecord{}
	scg.buylist = mtgban.BuylistRecord{}
	scg.client = NewSCGClient(guid, bearer)
	scg.client.SealedMode = true
	scg.MaxConcurrency = defaultConcurrency
	return &scg
}

func (scg *StarcitygamesSealed) printf(format string, a ...interface{}) {
	if scg.LogCallback != nil {
		scg.LogCallback("[SCGSealed] "+format, a...)
	}
}

func (scg *StarcitygamesSealed) processPage(channel chan<- responseChan, page int) error {
	results, err := scg.client.GetPage(page)
	if err != nil {
		return err
	}

	for _, result := range results {
		if len(result.Document.ProductType) == 0 {
			return errors.New("malformed product_type")
		}
		if result.Document.ProductType[0] == "Singles" {
			scg.printf("Skipping product_type %s", result.Document.ProductType[0])
			continue
		}

		if len(result.Document.ItemDisplayName) == 0 {
			return errors.New("malformed item_display_name")
		}
		if len(result.Document.UniqueID) == 0 {
			return errors.New("malformed unique_id")
		}

		if len(result.Document.URLDetail) == 0 {
			return errors.New("malformed url_detail")
		}

		productName := result.Document.ItemDisplayName[0]
		id := result.Document.UniqueID[0]
		urlPath := result.Document.URLDetail[0]

		if !strings.Contains(urlPath, "-mtg-") {
			continue
		}

		uuid, err := preprocessSealed(productName)
		if err != nil {
			continue
		}

		if uuid == "" {
			scg.printf("unable to parse %s", productName)
			continue
		}

		for _, attribute := range result.Document.HawkChildAttributes {
			if len(attribute.Price) == 0 {
				return errors.New("malformed price")
			}
			if len(attribute.Qty) == 0 {
				return errors.New("malformed qty")
			}
			priceStr := attribute.Price[0]
			qtyStr := attribute.Qty[0]

			price, err := mtgmatcher.ParsePrice(priceStr)
			if err != nil {
				scg.printf("invalid price for %s: %s", productName, err.Error())
				continue
			}

			qty, err := strconv.Atoi(qtyStr)
			if err != nil {
				scg.printf("invalid price for %s: %s", productName, err.Error())
				continue
			}

			if qty == 0 || price == 0 {
				continue
			}

			link := SCGProductURL(result.Document.URLDetail, attribute.VariantSKU, scg.Affiliate)

			out := responseChan{
				cardId: uuid,
				invEntry: &mtgban.InventoryEntry{
					Price:      price,
					Quantity:   qty,
					OriginalId: id,
					URL:        link,
				},
			}
			channel <- out
		}
	}

	return nil
}

func (scg *StarcitygamesSealed) scrape() error {
	totalPages, err := scg.client.NumberOfPages()
	if err != nil {
		return err
	}
	scg.printf("Found %d pages", totalPages)

	pages := make(chan int)
	results := make(chan responseChan)
	var wg sync.WaitGroup

	for i := 0; i < scg.MaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			for page := range pages {
				scg.printf("Processing page %d", page)
				err := scg.processPage(results, page)
				if err != nil {
					scg.printf("%v", err)
				}
			}
			wg.Done()
		}()
	}

	go func() {
		for i := 1; i <= totalPages; i++ {
			pages <- i
		}
		close(pages)

		wg.Wait()
		close(results)
	}()

	for record := range results {
		err := scg.inventory.Add(record.cardId, record.invEntry)
		if err != nil && !record.ignoreErr {
			scg.printf("%s", err.Error())
		}
	}

	scg.inventoryDate = time.Now()

	return nil
}

func (scg *StarcitygamesSealed) Inventory() (mtgban.InventoryRecord, error) {
	if len(scg.inventory) > 0 {
		return scg.inventory, nil
	}

	err := scg.scrape()
	if err != nil {
		return nil, err
	}

	return scg.inventory, nil

}

func (scg *StarcitygamesSealed) processBLPage(channel chan<- responseChan, page int) error {
	search, err := scg.client.SearchAll(page, defaultRequestLimit)
	if err != nil {
		return err
	}

	for _, hit := range search.Hits {
		productName := hit.Name

		uuid, err := preprocessSealed(productName)
		if err != nil {
			continue
		}

		if uuid == "" {
			scg.printf("unable to parse %s", productName)
			continue
		}

		link, _ := url.JoinPath(
			buylistBookmark,
			url.QueryEscape(hit.Name),
			",/0/0/0", // various faucets (hot list, rarity, bulk etc)
			fmt.Sprint(hit.SetID),
			",",           // unclear
			hit.Language,  // language ofc<D-x>
			"0/999999.99", // min/max price range
			",",           // finish
			"default",
		)

		for _, result := range hit.Variants {
			var priceRatio, sellPrice float64
			price := result.BuyPrice
			trade := result.TradePrice

			invCards := scg.inventory[uuid]
			for _, invCard := range invCards {
				sellPrice = invCard.Price
				break
			}
			if sellPrice > 0 {
				priceRatio = price / sellPrice * 100
			}

			channel <- responseChan{
				cardId: uuid,
				buyEntry: &mtgban.BuylistEntry{
					BuyPrice:   price,
					TradePrice: trade,
					Quantity:   0,
					PriceRatio: priceRatio,
					URL:        link,
				},
			}
		}
	}
	return nil
}

func (scg *StarcitygamesSealed) parseBL() error {
	search, err := scg.client.SearchAll(0, 1)
	if err != nil {
		return err
	}
	scg.printf("Parsing %d products", search.EstimatedTotalHits)

	pages := make(chan int)
	results := make(chan responseChan)
	var wg sync.WaitGroup

	for i := 0; i < scg.MaxConcurrency; i++ {
		wg.Add(1)
		go func() {
			for page := range pages {
				scg.printf("Processing page %d", page)
				err := scg.processBLPage(results, page)
				if err != nil {
					scg.printf("%v", err)
				}
			}
			wg.Done()
		}()
	}

	go func() {
		for j := 0; j < search.EstimatedTotalHits; j += defaultRequestLimit {
			pages <- j
		}
		close(pages)

		wg.Wait()
		close(results)
	}()

	for record := range results {
		err := scg.buylist.Add(record.cardId, record.buyEntry)
		if err != nil {
			scg.printf("%s", err.Error())
			continue
		}
	}

	scg.buylistDate = time.Now()

	return nil
}

func (scg *StarcitygamesSealed) Buylist() (mtgban.BuylistRecord, error) {
	if len(scg.buylist) > 0 {
		return scg.buylist, nil
	}

	err := scg.parseBL()
	if err != nil {
		return nil, err
	}

	return scg.buylist, nil
}

func (scg *StarcitygamesSealed) Info() (info mtgban.ScraperInfo) {
	info.Name = "Star City Games"
	info.Shorthand = "SCGSealed"
	info.InventoryTimestamp = &scg.inventoryDate
	info.BuylistTimestamp = &scg.buylistDate
	info.SealedMode = true
	return
}
