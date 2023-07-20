package main

import (
	"bufio"
	"encoding/xml"
	"example.com/m/v2/internal/xmls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getExchangeRate(rawCode, rawDate string) (string, error) {
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " ")
		switch data[0] {
		case "currency_rates":
			result, err := getExchangeRate(data[1], data[2])
			if err != nil {
				fmt.Printf("failed to get exchange rate: %v\n", err)
				continue
			}
			fmt.Println(result)
		case "close":
			return
		default:
			fmt.Println("Wrong command!")
		}
	}
}
