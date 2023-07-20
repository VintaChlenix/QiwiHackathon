package app

import (
	"encoding/xml"
	"example.com/m/v2/internal/xmls"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetExchangeRate(rawCode, rawDate string) (string, error) {
	_, code, _ := strings.Cut(rawCode, "=")
	_, date, _ := strings.Cut(rawDate, "=")
	splittedDate := strings.Split(date, "-")
	url := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily.asp?date_req=%s/%s/%s", splittedDate[2], splittedDate[1], splittedDate[0])
	valCurs, err := getValCursFromURL(url)
	if err != nil {
		return "", fmt.Errorf("failed to get exchange schema: %v", err)
	}

	for _, elem := range valCurs.Valute {
		if elem.CharCode == code {
			return elem.Value, nil
		}
	}

	return "", fmt.Errorf("failed to find your currency")
}

func getValCursFromURL(url string) (*xmls.ValCurs, error) {
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
