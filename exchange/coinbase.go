package exchange

type coinbaseService struct{
    Alt  []string
    Base []string
}

type coinbaseTicker struct {
    Data struct {
        Base       string `json:"base"`
        Currency   string `json:"currency"`
        Last       float64 `json:"amount,string"`
    } `json:"data"`
}

// Pairs supported by the exchange: quote -> base for easier maintainence.
var coinbaseCurr = map[string][]string {
    AUD: { BTC, ETH },
    USD: { BCH, BTC, ETH },
}

func (s coinbaseService) getPairs() (error, map[string][]string) {
    return nil, coinbaseCurr
}

func (s coinbaseService) getLastPrice(base string, quot string) (error, float64) {
    var custom coinbaseTicker
    urlP := base + "-" + quot
    err := GetJson("https://api.coinbase.com/v2/prices/" + urlP + "/spot", &custom)

    if (err != nil) {
        return err, -1
    }
    return err, custom.Data.Last
}
