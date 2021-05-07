package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hendychua/finantier-assessment/data-service/alphavantage"
	"github.com/hendychua/finantier-assessment/data-service/encrypt"
)

var alphaVantageClient alphavantage.AlphaVantageClient

var encryptionClient encrypt.EncryptionClient

// StockData is struct that holds a subset of data from mutiple structs in one struct.
type StockData struct {
	CurrentValue         string `json:"price"`
	Variation            string `json:"variation"`
	PrevClose            string `json:"prevClose"`
	Open                 string `json:"open"`
	Volume               string `json:"volume"`
	MarketCapitalization string `json:"marketCapitalization"`
}

func getSymbol(c *gin.Context) {
	symbol := c.Params.ByName("symbol")

	globalQuote, err := alphaVantageClient.GetGlobalQuote(symbol)
	if err != nil {
		log.Printf("Error while getting quote for '%s': %s\n", symbol, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if globalQuote.GlobalQuote.Symbol == "" {
		c.String(http.StatusNotFound, fmt.Sprintf("No such symbol found: '%s'", symbol))
		return
	}

	overview, err := alphaVantageClient.GetOverview(symbol)
	if err != nil {
		log.Printf("Error while getting overview for '%s': %s\n", symbol, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	stockData := StockData{
		CurrentValue:         globalQuote.GlobalQuote.Price,
		Variation:            globalQuote.GlobalQuote.ChangePercent,
		PrevClose:            globalQuote.GlobalQuote.PrevClose,
		Open:                 globalQuote.GlobalQuote.Open,
		Volume:               globalQuote.GlobalQuote.Volume,
		MarketCapitalization: overview.MarketCapitalization,
	}

	log.Printf("Stock data for symbol '%s': '%+v'\n", symbol, stockData)

	bytes, err := json.Marshal(stockData)
	if err != nil {
		log.Printf("Error while massaging stockData: '%+v': %s\n", stockData, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	payload := string(bytes)
	encryptedPayload, err := encryptionClient.Encrypt(&payload)
	if err != nil {
		log.Printf("Error while encrypting payload: '%s': %s\n", payload, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.String(http.StatusOK, *encryptedPayload)
}

// setupRouter sets up the API endpoints.
func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/symbol/:symbol", getSymbol)

	return r
}

func main() {
	key, found := os.LookupEnv("ALPHA_VANTAGE_API_KEY")
	if !found {
		log.Fatalln("Missing ALPHA_VANTAGE_API_KEY in environment.")
	}

	apiServer, found := os.LookupEnv("ENCRYPTION_SERVICE_HOST")
	if !found {
		log.Fatalln("Missing ENCRYPTION_SERVICE_HOST in environment.")
	}

	skipSSLVerify := os.Getenv("SKIP_SSL_VERIFY") == "1"
	log.Printf("Skip SSL verify: %t\n", skipSSLVerify)

	// Only for convenience
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSLVerify},
	}

	alphaVantageClient = alphavantage.AlphaVantageClient{
		APIKey: key,
		Client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: tr,
		},
	}

	encryptionClient = encrypt.EncryptionClient{
		APIServer: apiServer,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	r := setupRouter()
	r.Run(":8081")
}
