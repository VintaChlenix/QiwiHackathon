package main

import (
	"bufio"
	"example.com/m/v2/internal/app"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " ")
		switch data[0] {
		case "currency_rates":
			result, err := app.GetExchangeRate(data[1], data[2])
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
