package encrypt

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/assert"
)

// MockClient is a client that "encrypts" data by simply appending "_encrypted" to the original input.
type MockClient struct{}

func (c MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	content := string(bytes)
	encrypted := content + "_encrypted"
	response := &http.Response{
		Status:     "200 0K",
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(encrypted)),
	}

	return response, nil
}

func TestEncrypt(t *testing.T) {
	client := EncryptionClient{
		APIServer: "http://localhost:8080/",
		Client:    &MockClient{},
	}

	payload := "test"
	result, err := client.Encrypt(&payload)
	assert.Equal(t, nil, err)
	assert.Equal(t, "test_encrypted", *result)
}
