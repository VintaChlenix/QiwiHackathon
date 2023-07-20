package app

import (
	"encoding/xml"
	"example.com/m/v2/internal/xmls"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetExchangeRate(data []string) (string, error) {
	_, code, _ := strings.Cut(data[1], "=")
	_, date, _ := strings.Cut(data[2], "=")
	splittedDate := strings.Split(date, "-")
	url := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily.asp?date_req=%s/%s/%s", splittedDate[2], splittedDate[1], splittedDate[0])

	curRates, err := getCurrenciesRatesFromURL(url)
	if err != nil {
		return "", fmt.Errorf("failed to get currencies rates: %v", err)
	}

	for _, elem := range curRates.Valute {
		if elem.CharCode == code {
			return elem.Value, nil
		}
	}

	return "", fmt.Errorf("failed to find your currency")
}

func getCurrenciesRatesFromURL(url string) (*xmls.ValCurs, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get xml scheme: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error: %v", response.StatusCode)
	}

	schema, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}

	exchangeRates := new(xmls.ValCurs)
	if err := xml.Unmarshal(schema, &exchangeRates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml: %w", err)
	}

	return exchangeRates, nil
}
