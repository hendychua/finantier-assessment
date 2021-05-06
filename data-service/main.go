package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hendychua/finantier-assessment/data-service/alphavantage"
)

// StockData is struct that holds a subset of data from mutiple structs in one struct.
type StockData struct {
	CurrentValue         string `json:"price"`
	Variation            string `json:"variation"`
	PrevClose            string `json:"prevClose"`
	Open                 string `json:"open"`
	Volume               string `json:"volume"`
	MarketCapitalization string `json:"marketCapitalization"`
}

func main() {
	key, found := os.LookupEnv("ALPHA_VANTAGE_API_KEY")
	if !found {
		log.Fatalln("Missing ALPHA_VANTAGE_API_KEY in environment.")
	}

	alphaVantageClient := alphavantage.AlphaVantageClient{
		APIKey: key,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	globalQuote, err := alphaVantageClient.GetGlobalQuote("TSLA")
	if err != nil {
		panic(err)
	}

	overview, err := alphaVantageClient.GetOverview("TSLA")
	if err != nil {
		panic(err)
	}

	stockData := StockData{
		CurrentValue:         globalQuote.GlobalQuote.Price,
		Variation:            globalQuote.GlobalQuote.ChangePercent,
		PrevClose:            globalQuote.GlobalQuote.PrevClose,
		Open:                 globalQuote.GlobalQuote.Open,
		Volume:               globalQuote.GlobalQuote.Volume,
		MarketCapitalization: overview.MarketCapitalization,
	}

	fmt.Printf("%+v\n%+v\n%+v\n", *globalQuote, *overview, stockData)
}
