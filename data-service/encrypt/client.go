package encrypt

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type IClient interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type EncryptionClient struct {
	APIServer string
	Client    IClient
}

func (c EncryptionClient) Encrypt(payload *string) (*string, error) {
	u, err := url.Parse(c.APIServer)
	if err != nil {
		return nil, err
	}

	u.Path = "/encrypt"
	url := u.String()

	resp, err := c.Client.Post(url, "application/json", strings.NewReader(*payload))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error from encrypting service while encrypting payload '%s': %s", *payload, resp.Status)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	encrypted := string(bytes)
	return &encrypted, nil

}
