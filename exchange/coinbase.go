package exchange

type coinbaseService struct{}

type coinbaseTicker struct {
    Data struct {
        Base       string `json:"base"`
        Currency   string `json:"currency"`
        Last       float64 `json:"amount,string"`
    } `json:"data"`
}

func (s coinbaseService) getLastPrice(pair TradingPair) (error, float64) {
    var custom coinbaseTicker
    urlP := pair.Pair("-")
    err := GetJson("https://api.coinbase.com/v2/prices/" + urlP + "/spot", &custom)

    if (err != nil) {
        return err, -1
    }
    return err, custom.Data.Last
}
