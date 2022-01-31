package main

import (
	"fmt"
	"os"

	"pow/pkg/client"
)

func main() {
	quoteHost := os.Getenv("QUOTE_HOST")

	cl := client.NewQuoteClient(quoteHost)
	quote, err := cl.GetQuote()
	if err != nil {
		fmt.Println("got error", err)
		return
	}
	fmt.Println(quote)
}
