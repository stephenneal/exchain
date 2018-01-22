package fiat

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

var (
	current  string
	err      error
	rates    Rates
	response *http.Response
	body     []byte
)

func GetRates() {
	current = "USD"

	// Use api.fixer.io to get a JSON response
	response, err = http.Get("http://api.fixer.io/latest?base=" + current)
	if err != nil {
		rlog.Error(err)
	}
	defer response.Body.Close()

	// Read the data into a byte slice(string)
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		rlog.Error(err)
	}

	// Unmarshal the JSON byte slice to a predefined struct
	err = json.Unmarshal(body, &rates)
	if err != nil {
		rlog.Error(err)
	}

	// Everything accessible in struct now
	rlog.Info("\n==== Currency Rates ====\n")

	rlog.Info("Base: \t", rates.Base)
	rlog.Info("Date: \t", rates.Date)
	rlog.Info("USD:  \t", rates.Currencies.USD)
	rlog.Info("AUD:  \t", rates.Currencies.AUD)
	rlog.Info("CAD:  \t", rates.Currencies.CAD)
	rlog.Info("CHF:  \t", rates.Currencies.CHF)
	rlog.Info("EUR:  \t", rates.Currencies.EUR)
	rlog.Info("RUB:  \t", rates.Currencies.RUB)
	rlog.Info("JPY:  \t", rates.Currencies.JPY)
	rlog.Info("NZD:  \t", rates.Currencies.NZD)
}