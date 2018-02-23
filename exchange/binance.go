package exchange

type binanceService struct{
    Alt  []string
    Base []string
}

type binanceTicker struct {
    Symbol    string `json:"symbol"`
    Last      float64 `json:"price,string"`
}

// Pairs supported by the exchange: quote -> base for easier maintainence.
var binanceCurr = map[string][]string {
    BTC  : { ETH },
    USDT : { BTC, ETH },
}

func (s binanceService) getPairs() (error, map[string][]string) {
    return nil, binanceCurr
}

func (s binanceService) getLastPrice(base string, quot string) (error, float64) {
    var custom binanceTicker
    urlP := base + quot
    err := GetJson("https://api.binance.com/api/v3/ticker/price?symbol=" + urlP, &custom)

    if (err != nil) {
    	return err, -1
    }
    return err, custom.Last
}
