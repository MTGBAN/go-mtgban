package starcitygames

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	scgInventoryURL = "https://essearchapi-na.hawksearch.com/api/v2/search"
	scgBuylistURL   = "https://search.starcitygames.com/indexes/sell_list_products/search"

	maxResultsPerPage = 96
)

type SCGClient struct {
	client *http.Client
	guid   string
	bearer string
}

func NewSCGClient(guid, bearer string) *SCGClient {
	scg := SCGClient{}
	scg.client = cleanhttp.DefaultClient()
	scg.guid = guid
	scg.bearer = bearer
	return &scg
}

// https://bridgeline.atlassian.net/wiki/spaces/HSKB/pages/3462479664/Hawksearch+v4.0+-+Search+API
type scgRetailRequest struct {
	Keyword         string              `json:"Keyword"`
	FacetSelections map[string][]string `json:"FacetSelections"`
	PageNo          int                 `json:"PageNo"`
	MaxPerPage      int                 `json:"MaxPerPage"`
	ClientGUID      string              `json:"clientguid"`
}

type scgRetailResponse struct {
	Pagination struct {
		NofResults  int `json:"NofResults"`
		CurrentPage int `json:"CurrentPage"`
		MaxPerPage  int `json:"MaxPerPage"`
		NofPages    int `json:"NofPages"`
	} `json:"Pagination"`
	Results []scgRetailResult `json:"Results"`
}

type scgRetailResult struct {
	Document struct {
		Subtitle            []string `json:"subtitle"`
		UniqueID            []string `json:"unique_id"`
		CardName            []string `json:"card_name"`
		Language            []string `json:"language"`
		Set                 []string `json:"set"`
		CollectorNumber     []string `json:"collector_number"`
		Finish              []string `json:"finish"`
		ProductType         []string `json:"product_type"`
		URLDetail           []string `json:"url_detail"`
		HawkChildAttributes []struct {
			Price           []string `json:"price"`
			ProdID          []string `json:"prod_id"`
			VariantSku      []string `json:"variant_sku"`
			Qty             []string `json:"qty"`
			VariantLanguage []string `json:"variant_language"`
			Condition       []string `json:"condition"`
		} `json:"hawk_child_attributes"`
	} `json:"Document"`
	IsVisible bool `json:"IsVisible"`
}

func (scg *SCGClient) sendRetailRequest(page int) (*scgRetailResponse, error) {
	q := scgRetailRequest{
		ClientGUID: scg.guid,
		MaxPerPage: maxResultsPerPage,
		PageNo:     page,
		FacetSelections: map[string][]string{
			"variant_instockonly": {"Yes"},
			"product_type":        {"Singles"},
			"game":                {"Magic: The Gathering"},
		},
	}

	payload, err := json.Marshal(&q)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, scgInventoryURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-HawkSearch-IgnoreTracking", "true")
	req.Header.Add("Content-Type", "application/json")

	resp, err := scg.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var search scgRetailResponse
	err = json.Unmarshal(data, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}

func (scg *SCGClient) NumberOfPages() (int, error) {
	response, err := scg.sendRetailRequest(0)
	if err != nil {
		return 0, err
	}
	return response.Pagination.NofPages, nil
}

func (scg *SCGClient) GetPage(page int) ([]scgRetailResult, error) {
	response, err := scg.sendRetailRequest(page)
	if err != nil {
		return nil, err
	}
	return response.Results, nil
}

type SCGSearchRequest struct {
	Q                string   `json:"q"`
	Filter           string   `json:"filter"`
	MatchingStrategy string   `json:"matchingStrategy"`
	Limit            int      `json:"limit"`
	Offset           int      `json:"offset"`
	Sort             []string `json:"sort"`
}

type SCGSearchResponse struct {
	Message            string    `json:"message,omitempty"`
	Code               string    `json:"code,omitempty"`
	Type               string    `json:"type,omitempty"`
	Link               string    `json:"link,omitempty"`
	Hits               []SCGCard `json:"hits"`
	Query              string    `json:"query"`
	ProcessingTimeMs   int       `json:"processingTimeMs"`
	Limit              int       `json:"limit"`
	Offset             int       `json:"offset"`
	EstimatedTotalHits int       `json:"estimatedTotalHits"`
}

type SCGCard struct {
	Name            string           `json:"name"`
	ID              int              `json:"id"`
	Subtitle        string           `json:"subtitle"`
	Sku             string           `json:"sku"`
	ProductType     string           `json:"product_type"`
	CardName        string           `json:"card_name"`
	Finish          string           `json:"finish"`
	Language        string           `json:"language"`
	CollectorNumber string           `json:"collector_number"`
	Rarity          string           `json:"rarity"`
	SetID           int              `json:"set_id"`
	SetName         string           `json:"set_name"`
	SetReleaseDate  int              `json:"set_release_date"`
	SetSymbol       string           `json:"set_symbol"`
	IsBuying        int              `json:"is_buying"`
	Hotlist         int              `json:"hotlist"`
	Variants        []SCGCardVariant `json:"variants"`
}

type SCGCardVariant struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Subtitle           string  `json:"subtitle"`
	VariantName        string  `json:"variant_name"`
	VariantValue       string  `json:"variant_value"`
	Sku                string  `json:"sku"`
	IsBuying           int     `json:"is_buying"`
	Hotlist            float64 `json:"hotlist"`
	BuyPrice           float64 `json:"buy_price"`
	TradePrice         float64 `json:"trade_price"`
	BonusCalculationID int     `json:"bonus_calculation_id"`
}

func (scg *SCGClient) SearchAll(offset, limit int) (*SCGSearchResponse, error) {
	q := SCGSearchRequest{
		Filter:           "is_buying = 1 AND (product_type = \"Singles\") AND ((language = \"en\") OR (language = \"ja\"))",
		MatchingStrategy: "all",
		Limit:            limit,
		Offset:           offset,
		Sort:             []string{"name:asc", "set_name:asc", "finish:desc"},
	}
	payload, err := json.Marshal(&q)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, scgBuylistURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+scg.bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err := scg.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var search SCGSearchResponse
	err = json.Unmarshal(data, &search)
	if err != nil {
		return nil, err
	}

	if search.Message != "" {
		return nil, fmt.Errorf(search.Message)
	}

	return &search, nil
}
