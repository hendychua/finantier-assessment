package alphavantage

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IClient interface {
	Get(url string) (resp *http.Response, err error)
}

type GlobalQuote struct {
	Symbol        string `json:"01. symbol"`
	Open          string `json:"02. open"`
	High          string `json:"03. high"`
	Low           string `json:"04. low"`
	Price         string `json:"05. price"`
	Volume        string `json:"06. volume"`
	Day           string `json:"07. latest trading day"`
	PrevClose     string `json:"08. previous close"`
	Change        string `json:"09. change"`
	ChangePercent string `json:"10. change percent"`
}

type AVGlobalQuoteResponse struct {
	GlobalQuote GlobalQuote `json:"Global Quote"`
}

type AVOverviewResponse struct {
	Symbol               string `json:"Symbol"`
	MarketCapitalization string `json:"MarketCapitalization"`
}

const avApiFormat = "https://www.alphavantage.co/query?function=%s&symbol=%s&apikey=%s"

type AlphaVantageClient struct {
	APIKey string
	Client IClient
}

func (c AlphaVantageClient) GetGlobalQuote(symbol string) (*AVGlobalQuoteResponse, error) {
	url := fmt.Sprintf(avApiFormat, "GLOBAL_QUOTE", symbol, c.APIKey)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	avGlobalQuoteResponse := AVGlobalQuoteResponse{}
	err = json.NewDecoder(resp.Body).Decode(&avGlobalQuoteResponse)
	if err != nil {
		return nil, err
	}

	return &avGlobalQuoteResponse, nil
}

func (c AlphaVantageClient) GetOverview(symbol string) (*AVOverviewResponse, error) {
	url := fmt.Sprintf(avApiFormat, "OVERVIEW", symbol, c.APIKey)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	avOverviewResponse := AVOverviewResponse{}
	err = json.NewDecoder(resp.Body).Decode(&avOverviewResponse)
	if err != nil {
		return nil, err
	}

	return &avOverviewResponse, nil
}
