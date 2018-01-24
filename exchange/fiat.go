package exchange

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

    "github.com/romana/rlog"
)

type Rates struct {
	Base       string `json:"base"`
	Date       string `json:"date"`
	Currencies struct {
		AUD float64 `json:"AUD"`
		CAD float64 `json:"CAD"`
		CHF float64 `json:"CHF"`
		EUR float64 `json:"EUR"`
		NZD float64 `json:"NZD"`
		RUB float64 `json:"RUB"`
		JPY float64 `json:"JPY"`
		USD float64 `json:"USD"`
	} `json:"rates"`
}

func GetFiatRates(base string) (error, Rates) {
	var err      error
	var rates    Rates
	var response *http.Response
	var body     []byte

	// Use api.fixer.io to get a JSON response
	response, err = http.Get("http://api.fixer.io/latest?base=" + base)
	if err != nil {
		//rlog.Error(err)
	    return err, rates
	}
	defer response.Body.Close()

	// Read the data into a byte slice(string)
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		//rlog.Error(err)
	    return err, rates
	}

	// Unmarshal the JSON byte slice to a predefined struct
	err = json.Unmarshal(body, &rates)
	if err != nil {
		//rlog.Error(err)
	    return err, rates
	}

	// Everything accessible in struct now
	rlog.Debug("==== Currency Rates ====\n")

	rlog.Debug("Base: \t", rates.Base)
	rlog.Debug("Date: \t", rates.Date)
	rlog.Debug("USD:  \t", rates.Currencies.USD)
	rlog.Debug("AUD:  \t", rates.Currencies.AUD)
	rlog.Debug("CAD:  \t", rates.Currencies.CAD)
	rlog.Debug("CHF:  \t", rates.Currencies.CHF)
	rlog.Debug("EUR:  \t", rates.Currencies.EUR)
	rlog.Debug("RUB:  \t", rates.Currencies.RUB)
	rlog.Debug("JPY:  \t", rates.Currencies.JPY)
	rlog.Debug("NZD:  \t", rates.Currencies.NZD)
	return err, rates
}