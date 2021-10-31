package transaction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Currency struct {
	Success bool               `json:"success"`
	Base    string             `json:"source"`
	Rates   map[string]float64 `json:"quotes"`
}

func getCurrencyFromRub(needCurrency string) (float64, error) {
	searcherParams := url.Values{}
	searcherParams.Add("access_key", "b68d56ac99f42a02e979d0d708e2c3a5")

	resp, err := http.Get("http://api.currencylayer.com/live" + "?" + searcherParams.Encode())
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var curr Currency
	curr.Rates = make(map[string]float64, 100)
	err = json.Unmarshal(body, &curr)
	if err != nil {
		return 0, err
	}

	if !curr.Success {
		return 0, fmt.Errorf("no currency")
	}

	base := curr.Base
	need, ok := curr.Rates[base+needCurrency]
	if !ok {
		return 0, fmt.Errorf("no currency")
	}
	fmt.Println(need, ok)
	value := curr.Rates[base+"RUB"] / need

	return value, nil
}
