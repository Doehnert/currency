package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Currency struct {
	Code  string
	Name  string
	Rates map[string]float64
}

type MyCurrencyExchange struct {
	Currencies map[string]Currency
}

func NewMyCurrencyExchange() *MyCurrencyExchange {
	return &MyCurrencyExchange{
		Currencies: make(map[string]Currency),
	}
}

func (ce *MyCurrencyExchange) FetchAllCurrencies() error {
	resp, err := http.Get("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cs, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	csMap := make(map[string]string)
	err = json.Unmarshal(cs, &csMap)
	if err != nil {
		return err
	}

	i := 0
	for code, name := range csMap {
		if i > 100 {
			break
		}
		c := Currency{
			Code:  code,
			Name:  name,
			Rates: make(map[string]float64),
		}
		ce.Currencies[code] = c
		i++
	}
	return nil
}

func (c *Currency) FetchCurrencyRate() error {

	currencyCode := c.Code

	resp, err := http.Get(
		fmt.Sprintf("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/%s.json", currencyCode),
	)
	if err != nil {
		return err
	}

	rates, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ratesStruct := make(map[string]interface{})
	err = json.Unmarshal(rates, &ratesStruct)
	if err != nil {
		return err
	}

	ratesMap := make(map[string]float64)
	for code, rate := range ratesStruct[currencyCode].(map[string]interface{}) {
		ratesMap[code] = rate.(float64)
	}
	c.Rates = ratesMap
	return nil
}
