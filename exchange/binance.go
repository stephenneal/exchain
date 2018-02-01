package exchange

type binanceService struct{}

type binanceTicker struct {
    Symbol    string `json:"symbol"`
    Last      float64 `json:"price,string"`
}

func (s binanceService) getLastPrice(pair TradingPair) (error, float64) {
    var custom binanceTicker
    urlP := pair.Pair("")
    err := GetJson("https://api.binance.com/api/v3/ticker/price?symbol=" + urlP, &custom)

    if (err != nil) {
    	return err, -1
    }
    return err, custom.Last
}
