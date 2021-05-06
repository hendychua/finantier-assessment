package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var router *gin.Engine

// MockEncryptionClient is a mock client that returns whatever is used to initialise
// the struct. It is used for testing.
type MockEncryptionClient struct {
	EncryptionResult *string
	EncryptionError  error
	DecryptionResult *string
	DecryptionError  error
}

func (c MockEncryptionClient) Encrypt(plaintext *string) (*string, error) {
	return c.EncryptionResult, c.EncryptionError
}

func (c MockEncryptionClient) Decrypt(encrypted *string) (*string, error) {
	return c.DecryptionResult, c.DecryptionError
}

func init() {
	router = setupRouter()
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	return recorder
}

func TestEncryptPaylodError(t *testing.T) {
	// setup
	encryptionClient = MockEncryptionClient{
		nil, fmt.Errorf("Error encrypting"), nil, nil,
	}

	response := performRequest(router, "POST", "/encrypt", strings.NewReader(""))
	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestEncryptPaylod(t *testing.T) {
	// setup
	message := "{\"message\": \"Test secret message\"}"
	encryptionClient = MockEncryptionClient{
		&message, nil, nil, nil,
	}

	response := performRequest(router, "POST", "/encrypt", strings.NewReader(message))
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, message, response.Body.String())
}
