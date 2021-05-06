package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/hendychua/finantier-assessment/encryption-service/encrypt"
)

var encryptionClient encrypt.IEncryptionClient

// encryptPaylod encrypts the request body with AES256 and returns the encrypted value
func encryptPaylod(c *gin.Context) {
	bytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body := string(bytes)
	encrypted, err := encryptionClient.Encrypt(&body)
	if err != nil {
		log.Printf("Error happened during encryption: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.String(http.StatusOK, *encrypted)
}

// setupRouter sets up the API endpoints.
func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/encrypt", encryptPaylod)

	return r
}

func main() {
	key, found := os.LookupEnv("AES256_ENC_KEY")
	if !found {
		log.Fatalln("Missing a 32-byte encryption key. Specify it with environment variable 'AES256_ENC_KEY'")
	}
	encryptionClient = encrypt.AES256EncryptionClient{Key: []byte(key)}
	r := setupRouter()
	r.Run(":8080")
}
