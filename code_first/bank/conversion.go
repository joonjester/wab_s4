package bank

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var frankfurterAPI = "https://api.frankfurter.app"

type Currency string

const (
	EUR Currency = "eur"
	USD Currency = "usd"
	JPN Currency = "jpn"
	GBP Currency = "gbp"
)

type RatesResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
}

func ConvertCurrency(amount float64, base Currency, target Currency) (*float64, error) {
	url := fmt.Sprintf("%s/latest?amount=%f&from=%s&to=%s",
		frankfurterAPI, amount, base, target)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching rates: %s\n", err)

	}
	defer response.Body.Close()

	var rates RatesResponse
	if err := json.NewDecoder(response.Body).Decode(&rates); err != nil {
		return nil, fmt.Errorf("error decoding json: %s\n", err)
	}

	for _, convert := range rates.Rates {
		return &convert, nil
	}

	return nil, fmt.Errorf("something went wrong while converting")
}
