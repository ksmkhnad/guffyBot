package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// getExchangeRate fetches the exchange rate from the provided API.
func GetExchangeRate() (float64, error) {
	url := "https://api.bcc.kz:10443/bcc/production/v1/public/rates"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Add("authorization", "Bearer REPLACE_BEARER_TOKEN")
	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	// Parse the JSON response to extract the exchange rate
	var exchangeRateData map[string]interface{}
	err = json.Unmarshal(body, &exchangeRateData)
	if err != nil {
		return 0, err
	}

	// Assuming the response structure has a field named "exchange_rate"
	exchangeRate, ok := exchangeRateData["exchange_rate"].(float64)
	if !ok {
		return 0, fmt.Errorf("unable to parse exchange rate from API response")
	}

	return exchangeRate, nil
}
