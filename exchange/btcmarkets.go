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

func (s btcmService) name() string {
    return "BTC Markets"
}

func (s btcmService) getTicker(pair string) (error, Ticker) {
    var response btcmTicker
    err := GetJson("https://api.btcmarkets.net/market/" + pair + "/tick", &response)
    return err, response
}

func (t btcmTicker) LastPrice() float64 {
    return t.Last
}
