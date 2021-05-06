package alphavantage

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/assert"
)

type MockClient struct {
	Resp  *http.Response
	Error error
}

func (c *MockClient) Get(url string) (resp *http.Response, err error) {
	return c.Resp, c.Error
}

func TestGetGlobalQuote(t *testing.T) {
	avClient := AlphaVantageClient{
		APIKey: "xxxxxyyyyy",
		Client: &MockClient{
			Resp: &http.Response{
				Status:     "200 0K",
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(
					`{
                        "Global Quote": {
                            "01. symbol": "MSFT",
                            "05. price": "600.0"
                        }
                     }`,
				)),
			},
			Error: nil,
		},
	}

	expectedQuote := AVGlobalQuoteResponse{
		GlobalQuote: GlobalQuote{
			Symbol: "MSFT",
			Price:  "600.0",
		},
	}

	avGlobalQuote, err := avClient.GetGlobalQuote("MSFT")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedQuote, avGlobalQuote)
}

func TestGetOverview(t *testing.T) {
	avClient := AlphaVantageClient{
		APIKey: "xxxxxyyyyy",
		Client: &MockClient{
			Resp: &http.Response{
				Status:     "200 0K",
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(
					`{
                        "Symbol": "MSFT",
                        "MarketCapitalization": "654321"
                     }`,
				)),
			},
			Error: nil,
		},
	}

	expectedOverview := AVOverviewResponse{
		Symbol:               "MSFT",
		MarketCapitalization: "654321",
	}

	overview, err := avClient.GetOverview("MSFT")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedOverview, overview)
}
