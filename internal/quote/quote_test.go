package quote_test

import (
	"net/http"

	c "github.com/achannarasappa/ticker/internal/common"
	. "github.com/achannarasappa/ticker/internal/quote"
	. "github.com/achannarasappa/ticker/test/http"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	g "github.com/onsi/gomega/gstruct"
)

func mockResponseCurrencyGOOG() {
	response := `{
		"quoteResponse": {
			"result": [
				{
					"currency": "USD",
					"symbol": "GOOG"
				}
			],
			"error": null
		}
	}
	`
	responseURL := `https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&fields=regularMarketPrice,currency&symbols=GOOG`
	httpmock.RegisterResponder("GET", responseURL, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, response)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})
}

func mockResponseCurrencyEURUSD() {
	response := `{
		"quoteResponse": {
			"result": [
				{
					"quoteType": "CURRENCY",
					"quoteSourceName": "Delayed Quote",
					"currency": "EUR",
					"regularMarketPrice": 0.8891,
					"sourceInterval": 15,
					"exchangeDataDelayedBy": 0,
					"exchange": "CCY",
					"fullExchangeName": "CCY",
					"symbol": "USDEUR=X"
				}
			],
			"error": null
		}
	}
	
	`
	responseURL := `https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&fields=regularMarketPrice,currency&symbols=USDEUR=X`
	httpmock.RegisterResponder("GET", responseURL, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, response)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})
}

var _ = Describe("Quote", func() {

	var (
		dep c.Dependencies
	)

	BeforeEach(func() {
		dep = c.Dependencies{
			HttpClient: client,
		}
		MockResponseYahooQuotes()
	})

	Describe("GetAssetGroupQuote", func() {

		It("should get price quotes for each asset based on it's data source", func() {

			input := c.AssetGroup{
				SymbolsBySource: []c.AssetGroupSymbolsBySource{
					{
						Source: c.QuoteSourceYahoo,
						Symbols: []string{
							"GOOG",
							"RBLX",
						},
					},
					{
						Source: c.QuoteSourceUserDefined,
						Symbols: []string{
							"CASH",
							"PRIVATESHARES",
						},
					},
				},
			}
			output := GetAssetGroupQuote(dep)(input)

			Expect(output.AssetQuotes).To(HaveLen(2))

		})

	})

	Describe("GetAssetGroupsCurrencyRates", func() {

		It("should get currency conversion rates for each type of data source", func() {

			mockResponseCurrencyEURUSD()
			mockResponseCurrencyGOOG()
			input := []c.AssetGroup{
				{
					SymbolsBySource: []c.AssetGroupSymbolsBySource{
						{
							Source: c.QuoteSourceYahoo,
							Symbols: []string{
								"GOOG",
							},
						},
					},
				},
			}
			output, _ := GetAssetGroupsCurrencyRates(*client, input, "EUR")
			Expect(output).To(g.MatchAllKeys(g.Keys{
				"USD": g.MatchFields(g.IgnoreExtras, g.Fields{
					"FromCurrency": Equal("USD"),
					"ToCurrency":   Equal("EUR"),
					"Rate":         Equal(0.8891),
				}),
			}))

		})

	})

})
