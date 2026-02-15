package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Coiiap5e/photographer/internal/errors"
)

const cbrURL = "https://www.cbr-xml-daily.ru/daily_json.js"

type CBRResponse struct {
	Currencies map[string]Currency `json:"Valute"`
}

type Currency struct {
	ID       string  `json:"ID"`
	NumCode  string  `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Nominal  int     `json:"Nominal"`
	Name     string  `json:"Name"`
	Value    float64 `json:"Value"`
	Previous float64 `json:"Previous"`
}

// It returns the value of 1 USD in RUB.

func GetUSDRate(ctx context.Context) (float64, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", cbrURL, nil)
	if err != nil {
		return 0, errors.Wrap(err, errors.ErrCodeCurrencyAPIRequest, "failed to create HTTP request for CBR API")
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, errors.ErrCodeCurrencyAPIRequest, "failed to execute HTTP request to CBR API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("CBR API returned non-200 status code: %d", resp.StatusCode)
		return 0, errors.New(errors.ErrCodeCurrencyAPIResponse, msg)
	}

	var cbrResponse CBRResponse
	if err := json.NewDecoder(resp.Body).Decode(&cbrResponse); err != nil {
		return 0, errors.Wrap(err, errors.ErrCodeCurrencyAPIParsing, "failed to decode JSON response from CBR API")
	}

	usdCurrency, ok := cbrResponse.Currencies["USD"]
	if !ok {
		return 0, errors.New(errors.ErrCodeCurrencyNotFound, "USD currency not found in CBR API response")
	}

	if usdCurrency.Nominal == 0 {
		return 0, errors.New(errors.ErrCodeCurrencyAPIResponse, "received invalid nominal value for USD")
	}

	rate := usdCurrency.Value / float64(usdCurrency.Nominal)
	return rate, nil
}
