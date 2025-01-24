package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/scizorman/go-ndjson"
	"github.com/ulikunitz/xz"
	"golang.org/x/exp/maps"
	"google.golang.org/api/option"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mtgban/go-mtgban/abugames"
	"github.com/mtgban/go-mtgban/cardkingdom"
	"github.com/mtgban/go-mtgban/cardmarket"
	"github.com/mtgban/go-mtgban/cardsphere"
	"github.com/mtgban/go-mtgban/cardtrader"
	"github.com/mtgban/go-mtgban/coolstuffinc"
	"github.com/mtgban/go-mtgban/hareruya"
	"github.com/mtgban/go-mtgban/jupitergames"
	"github.com/mtgban/go-mtgban/magiccorner"
	"github.com/mtgban/go-mtgban/manapool"
	"github.com/mtgban/go-mtgban/mintcard"
	"github.com/mtgban/go-mtgban/mtgseattle"
	"github.com/mtgban/go-mtgban/mtgstocks"
	"github.com/mtgban/go-mtgban/ninetyfive"
	"github.com/mtgban/go-mtgban/sealedev"
	"github.com/mtgban/go-mtgban/starcitygames"
	"github.com/mtgban/go-mtgban/strikezone"
	"github.com/mtgban/go-mtgban/tcgplayer"
	"github.com/mtgban/go-mtgban/toamagic"
	"github.com/mtgban/go-mtgban/trollandtoad"
	"github.com/mtgban/go-mtgban/wizardscupboard"

	"github.com/mtgban/go-mtgban/mtgban"
	"github.com/mtgban/go-mtgban/mtgmatcher"
)

var date = time.Now().Format("2006-01-02")
var GCSBucket *storage.BucketHandle

var GlobalLogCallback mtgban.LogCallbackFunc

var MaxConcurrency = os.Getenv("MAX_CONCURRENCY")

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}()

type scraperOption struct {
	Enabled    bool
	OnlySeller bool
	OnlyVendor bool
	Init       func() (mtgban.Scraper, error)
}

var options = map[string]*scraperOption{
	"abugames": {
		Init: func() (mtgban.Scraper, error) {
			scraper := abugames.NewScraper()
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"cardkingdom": {
		Init: func() (mtgban.Scraper, error) {
			scraper := cardkingdom.NewScraper()
			scraper.LogCallback = GlobalLogCallback
			scraper.Partner = os.Getenv("CK_PARTNER")
			return scraper, nil
		},
	},
	"cardkingdom_sealed": {
		Init: func() (mtgban.Scraper, error) {
			scraper := cardkingdom.NewScraperSealed()
			scraper.LogCallback = GlobalLogCallback
			scraper.Partner = os.Getenv("CK_PARTNER")
			return scraper, nil
		},
	},
	"cardsphere": {
		Init: func() (mtgban.Scraper, error) {
			csphereToken := os.Getenv("CARDSPHERE_TOKEN")
			if csphereToken == "" {
				return nil, errors.New("missing CARDSPHERE_TOKEN env var")
			}
			scraper := cardsphere.NewScraper(csphereToken)
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"cardtrader": {
		Init: func() (mtgban.Scraper, error) {
			ctTokenBearer := os.Getenv("CARDTRADER_TOKEN_BEARER")
			if ctTokenBearer == "" {
				return nil, errors.New("missing CARDTRADER_TOKEN_BEARER env var")
			}

			scraper, err := cardtrader.NewScraperMarket(cardtrader.GameIdMagic, ctTokenBearer)
			if err != nil {
				return nil, err
			}
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"cardtrader_sealed": {
		Init: func() (mtgban.Scraper, error) {
			ctTokenBearer := os.Getenv("CARDTRADER_TOKEN_BEARER")
			if ctTokenBearer == "" {
				return nil, errors.New("missing CARDTRADER_TOKEN_BEARER env var")
			}

			scraper, err := cardtrader.NewScraperSealed(ctTokenBearer)
			if err != nil {
				return nil, err
			}
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"coolstuffinc": {
		OnlyVendor: true,
		Init: func() (mtgban.Scraper, error) {
			scraper := coolstuffinc.NewScraper(coolstuffinc.GameMagic)
			scraper.LogCallback = GlobalLogCallback
			scraper.Partner = os.Getenv("CSI_PARTNER")
			return scraper, nil
		},
	},
	"coolstuffinc_official": {
		Init: func() (mtgban.Scraper, error) {
			csiKey := os.Getenv("CSI_KEY")
			if csiKey == "" {
				return nil, errors.New("missing CSI_KEY env var")
			}

			scraper := coolstuffinc.NewScraperOfficial(csiKey)
			scraper.LogCallback = GlobalLogCallback
			scraper.Partner = os.Getenv("CSI_PARTNER")
			return scraper, nil
		},
	},
	"coolstuffinc_sealed": {
		OnlyVendor: true,
		Init: func() (mtgban.Scraper, error) {
			scraper := coolstuffinc.NewScraperSealed()
			scraper.LogCallback = GlobalLogCallback
			scraper.Partner = os.Getenv("CSI_PARTNER")
			return scraper, nil
		},
	},
	"hareruya": {
		Init: func() (mtgban.Scraper, error) {
			scraper, err := hareruya.NewScraper()
			if err != nil {
				return nil, err
			}
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"jupitergames": {
		OnlySeller: true,
		Init: func() (mtgban.Scraper, error) {
			scraper := jupitergames.NewScraper()
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"magiccorner": {
		Init: func() (mtgban.Scraper, error) {
			scraper, err := magiccorner.NewScraper()
			if err != nil {
				return nil, err
			}
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"manapool": {
		Init: func() (mtgban.Scraper, error) {
			scraper := manapool.NewScraper()
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"mintcard": {
		Init: func() (mtgban.Scraper, error) {
			scraper := mintcard.NewScraper()
			scraper.Partner = os.Getenv("MP_AFFILIATE")
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"mkm_index": {
		Init: func() (mtgban.Scraper, error) {
			mkmAppToken := os.Getenv("MKM_APP_TOKEN")
			mkmAppSecret := os.Getenv("MKM_APP_SECRET")
			if mkmAppToken == "" || mkmAppSecret == "" {
				return nil, errors.New("missing MKM_APP_TOKEN or MKM_APP_SECRET env vars")
			}

			scraper, err := cardmarket.NewScraperIndex(cardmarket.GameIdMagic, mkmAppToken, mkmAppSecret)
			if err != nil {
				return nil, err
			}
			scraper.Affiliate = "mtgban"
			scraper.LogCallback = GlobalLogCallback
			scraper.Affiliate = os.Getenv("MKM_AFFILIATE")
			return scraper, nil
		},
	},
	"mkm_sealed": {
		Init: func() (mtgban.Scraper, error) {
			mkmAppToken := os.Getenv("MKM_APP_TOKEN")
			mkmAppSecret := os.Getenv("MKM_APP_SECRET")
			if mkmAppToken == "" || mkmAppSecret == "" {
				return nil, errors.New("missing MKM_APP_TOKEN or MKM_APP_SECRET env vars")
			}

			scraper, err := cardmarket.NewScraperSealed(mkmAppToken, mkmAppSecret)
			if err != nil {
				return nil, err
			}
			scraper.LogCallback = GlobalLogCallback
			scraper.Affiliate = os.Getenv("MKM_AFFILIATE")
			return scraper, nil
		},
	},
	"mtgseattle": {
		Init: func() (mtgban.Scraper, error) {
			scraper := mtgseattle.NewScraper()
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"mtgstocks": {
		Init: func() (mtgban.Scraper, error) {
			scraper := mtgstocks.NewScraper()
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"ninetyfive": {
		Init: func() (mtgban.Scraper, error) {
			scraper, err := ninetyfive.NewScraper()
			if err != nil {
				return nil, err
			}
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"sealed_ev": {
		Init: func() (mtgban.Scraper, error) {
			banKey := os.Getenv("BAN_API_KEY")
			if banKey == "" {
				return nil, errors.New("missing BAN_API_KEY env var")
			}
			scraper := sealedev.NewScraper(banKey)
			scraper.Affiliate = os.Getenv("TCG_AFFILIATE")
			scraper.BuylistAffiliate = os.Getenv("CK_PARTNER")
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"starcitygames": {
		Init: func() (mtgban.Scraper, error) {
			scgGUID := os.Getenv("SCG_GUID")
			scgBearer := os.Getenv("SCG_BEARER")
			if scgGUID == "" || scgBearer == "" {
				return nil, errors.New("missing SCG_GUID or SCG_BEARER env var")
			}

			scraper := starcitygames.NewScraper(starcitygames.GameMagic, scgGUID, scgBearer)
			scraper.LogCallback = GlobalLogCallback
			scraper.Affiliate = os.Getenv("SCG_AFFILIATE")
			return scraper, nil
		},
	},
	"starcitygames_sealed": {
		Init: func() (mtgban.Scraper, error) {
			scgGUID := os.Getenv("SCG_GUID")
			scgBearer := os.Getenv("SCG_BEARER")
			if scgGUID == "" || scgBearer == "" {
				return nil, errors.New("missing SCG_GUID or SCG_BEARER env var")
			}

			scraper := starcitygames.NewScraperSealed(scgGUID, scgBearer)
			scraper.LogCallback = GlobalLogCallback
			scraper.Affiliate = os.Getenv("SCG_AFFILIATE")
			return scraper, nil
		},
	},
	"strikezone": {
		Init: func() (mtgban.Scraper, error) {
			scraper := strikezone.NewScraper(strikezone.GameMagic)
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"tcg_index": {
		Init: func() (mtgban.Scraper, error) {
			tcgPublicId := os.Getenv("TCGPLAYER_PUBLIC_ID")
			tcgPrivateId := os.Getenv("TCGPLAYER_PRIVATE_ID")
			if tcgPublicId == "" || tcgPrivateId == "" {
				return nil, errors.New("missing TCGPLAYER_PUBLIC_ID or TCGPLAYER_PRIVATE_ID env vars")
			}

			scraper := tcgplayer.NewScraperIndex(tcgPublicId, tcgPrivateId)
			scraper.LogCallback = GlobalLogCallback
			scraper.Affiliate = os.Getenv("TCG_AFFILIATE")
			num, _ := strconv.Atoi(MaxConcurrency)
			if num != 0 {
				scraper.MaxConcurrency = num
			}
			return scraper, nil
		},
	},
	"tcg_market": {
		Init: func() (mtgban.Scraper, error) {
			tcgPublicId := os.Getenv("TCGPLAYER_PUBLIC_ID")
			tcgPrivateId := os.Getenv("TCGPLAYER_PRIVATE_ID")
			mtgjsonTCGSKUFilepathBZ2 := os.Getenv("MTGJSON_TCGSKU_FILEPATH_BZ2")
			if tcgPublicId == "" || tcgPrivateId == "" || mtgjsonTCGSKUFilepathBZ2 == "" {
				return nil, errors.New("missing TCGPLAYER_PUBLIC_ID or TCGPLAYER_PRIVATE_ID or MTGJSON_TCGSKU_FILEPATH_BZ2 env vars")
			}

			scraper := tcgplayer.NewScraperMarket(tcgPublicId, tcgPrivateId)
			scraper.LogCallback = GlobalLogCallback
			scraper.Affiliate = os.Getenv("TCG_AFFILIATE")
			num, _ := strconv.Atoi(MaxConcurrency)
			if num != 0 {
				scraper.MaxConcurrency = num
			}
			var reader io.ReadCloser
			var err error
			if strings.HasPrefix(mtgjsonTCGSKUFilepathBZ2, "http") {
				resp, err := http.Get(mtgjsonTCGSKUFilepathBZ2)
				if err != nil {
					return nil, err
				}
				reader = resp.Body
			} else {
				reader, err = os.Open(mtgjsonTCGSKUFilepathBZ2)
				if err != nil {
					return nil, err
				}
			}
			defer reader.Close()
			skus, err := tcgplayer.LoadTCGSKUs(reader)
			if err != nil {
				return nil, err
			}
			scraper.SKUsData = skus.Data
			return scraper, nil
		},
	},
	"tcg_sealed": {
		Init: func() (mtgban.Scraper, error) {
			tcgPublicId := os.Getenv("TCGPLAYER_PUBLIC_ID")
			tcgPrivateId := os.Getenv("TCGPLAYER_PRIVATE_ID")
			if tcgPublicId == "" || tcgPrivateId == "" {
				return nil, errors.New("missing TCGPLAYER_PUBLIC_ID or TCGPLAYER_PRIVATE_ID env vars")
			}

			scraper := tcgplayer.NewScraperSealed(tcgPublicId, tcgPrivateId)
			scraper.LogCallback = GlobalLogCallback
			scraper.Affiliate = os.Getenv("TCG_AFFILIATE")
			return scraper, nil
		},
	},
	"tcg_syplist": {
		Init: func() (mtgban.Scraper, error) {
			tcgAuth := os.Getenv("TCGPLAYER_AUTH")
			if tcgAuth == "" {
				return nil, errors.New("missing TCGPLAYER_AUTH env var")
			}
			scraper := tcgplayer.NewScraperSYP(tcgAuth)
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"toamagic": {
		Init: func() (mtgban.Scraper, error) {
			scraper := toamagic.NewScraper()
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
	"trollandtoad": {
		Init: func() (mtgban.Scraper, error) {
			scraper := trollandtoad.NewScraper()
			scraper.LogCallback = GlobalLogCallback
			num, _ := strconv.Atoi(MaxConcurrency)
			if num != 0 {
				scraper.MaxConcurrency = num
			}
			return scraper, nil
		},
	},
	"wizardscupboard": {
		Init: func() (mtgban.Scraper, error) {
			scraper := wizardscupboard.NewScraper()
			scraper.LogCallback = GlobalLogCallback
			return scraper, nil
		},
	},
}

type inventoryElement struct {
	UUID string
	mtgban.InventoryEntry
}

type buylistElement struct {
	UUID string
	mtgban.BuylistEntry
}

func writeSellerToNDJSON(seller mtgban.Seller, w io.Writer) error {
	inventory, err := seller.Inventory()
	if err != nil {
		return err
	}

	var inventoryFlat []inventoryElement
	for uuid, entries := range inventory {
		for _, entry := range entries {
			inventoryFlat = append(inventoryFlat, inventoryElement{
				UUID:           uuid,
				InventoryEntry: entry,
			})
		}
	}

	output, err := ndjson.Marshal(inventoryFlat)
	if err != nil {
		return err
	}

	_, err = w.Write(output)
	return err
}

func writeVendorToNDJSON(vendor mtgban.Vendor, w io.Writer) error {
	buylist, err := vendor.Buylist()
	if err != nil {
		return err
	}

	var buylistFlat []buylistElement
	for uuid, entries := range buylist {
		for _, entry := range entries {
			buylistFlat = append(buylistFlat, buylistElement{
				UUID:         uuid,
				BuylistEntry: entry,
			})
		}
	}

	output, err := ndjson.Marshal(buylistFlat)
	if err != nil {
		return err
	}

	_, err = w.Write(output)
	return err
}

func dumpSeller(seller mtgban.Seller, outputPath, format string) error {
	fname := fmt.Sprintf("%s_%s_Inventory.%s", date, seller.Info().Shorthand, format)
	filePath := path.Join(outputPath, fname)

	var writer io.WriteCloser
	if GCSBucket == nil {
		log.Println("Dumping seller", seller.Info().Shorthand)
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		writer = file
	} else {
		log.Println("Uploading seller", seller.Info().Shorthand)
		writer = GCSBucket.Object(filePath).NewWriter(context.Background())
	}
	defer writer.Close()

	var err error
	switch format {
	case "json":
		err = mtgban.WriteSellerToJSON(seller, writer)
	case "csv":
		err = mtgban.WriteSellerToCSV(seller, writer)
	case "ndjson":
		err = writeSellerToNDJSON(seller, writer)
	}

	return err
}

func dumpVendor(vendor mtgban.Vendor, outputPath, format string) error {
	fname := fmt.Sprintf("%s_%s_Buylist.%s", date, vendor.Info().Shorthand, format)
	filePath := path.Join(outputPath, fname)

	var writer io.WriteCloser
	if GCSBucket == nil {
		log.Println("Dumping vendor", vendor.Info().Shorthand)
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		writer = file
	} else {
		log.Println("Uploading vendor", vendor.Info().Shorthand)
		writer = GCSBucket.Object(filePath).NewWriter(context.Background())
	}
	defer writer.Close()

	var err error
	switch format {
	case "json":
		err = mtgban.WriteVendorToJSON(vendor, writer)
	case "csv":
		err = mtgban.WriteVendorToCSV(vendor, writer)
	case "ndjson":
		err = writeVendorToNDJSON(vendor, writer)
	}

	return err
}

func dump(bc *mtgban.BanClient, outputPath, format string, meta bool) error {
	for _, seller := range bc.Sellers() {
		err := dumpSeller(seller, outputPath, format)
		if err != nil {
			return err
		}

		if meta && format != "json" {
			sellerMeta := mtgban.NewSellerFromInventory(nil, seller.Info())
			err := dumpSeller(sellerMeta, outputPath, "json")
			if err != nil {
				return err
			}
		}
	}

	for _, vendor := range bc.Vendors() {
		err := dumpVendor(vendor, outputPath, format)
		if err != nil {
			return err
		}

		if meta && format != "json" {
			vendorMeta := mtgban.NewVendorFromBuylist(nil, vendor.Info())
			err := dumpVendor(vendorMeta, outputPath, "json")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func run() int {
	for key, val := range options {
		flag.BoolVar(&val.Enabled, key, false, "Enable "+strings.Title(key))
	}

	mtgjsonOpt := flag.String("mtgjson", "", "Path to AllPrintings file")
	outputPathOpt := flag.String("output-path", "", "Path where to dump results")
	serviceAccOpt := flag.String("svc-acc", "", "Service account with write permission on the bucket")

	scrapersOpt := flag.String("scrapers", "", "Comma-separated list of scrapers to enable")
	sellersOpt := flag.String("sellers", "", "Comma-separated list of sellers to enable")
	vendorsOpt := flag.String("vendors", "", "Comma-separated list of vendors to enable")

	fileFormatOpt := flag.String("format", "json", "File format of the output files (json/csv/ndjson)")
	metaOpt := flag.Bool("meta", false, "When format is not json, output a second file for scraper metadata")

	devOpt := flag.Bool("dev", false, "Enable dev operations (debugging)")
	versionOpt := flag.Bool("v", false, "Print version information")
	listOpt := flag.Bool("l", false, "List all scrapers available")
	flag.Parse()

	log.Println("bantool version", Commit)
	if *versionOpt {
		return 0
	}

	if *listOpt {
		keys := maps.Keys(options)
		sort.Strings(keys)
		for _, key := range keys {
			var list []string
			if options[key].OnlyVendor {
				list = append(list, "❌")
			} else {
				list = append(list, "✅")
			}
			if options[key].OnlySeller {
				list = append(list, "❌")
			} else {
				list = append(list, "✅")
			}
			fmt.Println(list, key)
		}
		return 0
	}

	if *devOpt {
		GlobalLogCallback = log.Printf
	}

	switch *fileFormatOpt {
	case "json", "csv", "ndjson":
	default:
		log.Println("Invalid -format option, see -h for supported values")
		return 1
	}

	if *outputPathOpt == "" {
		log.Println("Missing output-path argument")
		return 1
	}

	u, err := url.Parse(*outputPathOpt)
	if err != nil {
		log.Println("cannot parse output-path", err)
		return 1
	}

	// If a service account file is passed in, create a bucket object and update the output path
	if *serviceAccOpt != "" {
		client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(*serviceAccOpt))
		if err != nil {
			log.Println("error creating the GCS client", err)
			return 1
		}

		if u.Scheme != "gs" {
			log.Println("unsupported scheme in output-path")
			return 1
		}

		GCSBucket = client.Bucket(u.Host)

		// Trim to avoid creating an empty directory in the bucket
		*outputPathOpt = strings.TrimPrefix(u.Path, "/")
	} else if u.Scheme != "" {
		log.Println("missing svc-acc file for cloud access")
		return 1
	} else {
		_, err := os.Stat(*outputPathOpt)
		if os.IsNotExist(err) {
			log.Println("output-path does not exist")
			return 1
		}
	}

	// Enable Scrapers or Sellers/Vendors
	scraps := strings.Split(*scrapersOpt, ",")
	for _, name := range scraps {
		if options[name] != nil {
			options[name].Enabled = true
		}
	}
	if *sellersOpt != "" {
		sells := strings.Split(*sellersOpt, ",")
		for _, name := range sells {
			if options[name] == nil {
				log.Fatalln("Seller", name, "not found")
			}
			options[name].Enabled = true
			options[name].OnlySeller = true
			options[name].OnlyVendor = false
		}
	}
	if *vendorsOpt != "" {
		vends := strings.Split(*vendorsOpt, ",")
		for _, name := range vends {
			if options[name] == nil {
				log.Fatalln("Vendor", name, "not found")
			}
			options[name].Enabled = true
			options[name].OnlySeller = false
			options[name].OnlyVendor = true
		}
	}

	bc := mtgban.NewClient()

	// Initialize the enabled scrapers
	for _, opt := range options {
		if opt.Enabled {
			scraper, err := opt.Init()
			if err != nil {
				log.Println(err)
				return 1
			}
			if opt.OnlySeller {
				bc.RegisterSeller(scraper)
			} else if opt.OnlyVendor {
				bc.RegisterVendor(scraper)
			} else {
				bc.Register(scraper)
			}
		}
	}

	if len(bc.Scrapers()) == 0 {
		log.Println("No scraper configured, run with -h for a list of commands")
		return 1
	}

	// Load static data
	if *mtgjsonOpt != "" {
		log.Println("Loading MTGJSON from", *mtgjsonOpt)
		err = mtgmatcher.LoadDatastoreFile(*mtgjsonOpt)
	} else {
		log.Println("Loading MTGJSON from net")
		err = loadMTGJSONfromNet()
	}
	if err != nil {
		log.Println("Couldn't load MTGJSON...")
		log.Println(err)
		return 1
	}

	// Load the data
	err = bc.Load()
	if err != nil {
		log.Println("Something didn't work while scraping...")
		log.Println(err)
		return 1
	}

	// Dump the results
	err = dump(bc, *outputPathOpt, *fileFormatOpt, *metaOpt)
	if err != nil {
		log.Println(err)
		return 1
	}

	log.Println("Completed")

	return 0
}

func main() {
	os.Exit(run())
}

const AllPrintingsURL = "https://mtgjson.com/api/v5/AllPrintings.json.xz"

func loadMTGJSONfromNet() error {
	resp, err := cleanhttp.DefaultClient().Get(AllPrintingsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader, err := xz.NewReader(resp.Body)
	if err != nil {
		return err
	}

	return mtgmatcher.LoadDatastore(reader)
}
