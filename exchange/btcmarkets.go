package exchange

type btcmService struct{}

type btcmTicker struct {
	BestBid    float64 `json:"bestBid"`
	BestAsk    float64 `json:"bestAsk"`
	Last       float64 `json:"lastPrice"`
	Currency   string  `json:"currency"`
	Instrument string  `json:"instrument"`
	Timestamp  int     `json:"timestamp"`
	Volume24H  float64 `json:"volume24h"`
}

// Pairs supported by the exchange: quote -> base for easier maintainence.
var btcmCurr = map[string][]string {
    AUD: { BCH, BTC, ETH },
}

func (s btcmService) getPairs() (error, map[string][]string) {
    return nil, btcmCurr
}

func (s btcmService) getLastPrice(base string, quot string) (error, float64) {
    var custom btcmTicker
    err := GetJson("https://api.btcmarkets.net/market/" + base + "/" + quot + "/tick", &custom)

    if (err != nil) {
    	return err, -1
    }
    return err, custom.Last
}
