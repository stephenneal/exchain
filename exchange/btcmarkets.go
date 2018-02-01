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

func (s btcmService) getLastPrice(pair TradingPair) (error, float64) {
    var custom btcmTicker
    err := GetJson("https://api.btcmarkets.net/market/" + pair.Pair("/") + "/tick", &custom)

    if (err != nil) {
    	return err, -1
    }
    return err, custom.Last
}

func (t btcmTicker) LastPrice() float64 {
    return t.Last
}
