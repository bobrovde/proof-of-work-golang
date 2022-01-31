package quote

import (
	"encoding/csv"
	"math/rand"
	"os"

	"pow/domain"
)

func NewCSVQuoter() (CSVQuoter, error) {
	file, err := os.Open("./pkg/quote/quotes.csv")
	if err != nil {
		return CSVQuoter{}, err
	}

	r := csv.NewReader(file)
	quotes, err := r.ReadAll()
	if err != nil {
		return CSVQuoter{}, err
	}

	return CSVQuoter{
		totalQuotes: len(quotes),
		quotes:      quotes,
	}, nil
}

type CSVQuoter struct {
	totalQuotes int
	quotes      [][]string
}

func (q CSVQuoter) GetRandomQuote() domain.Quote {
	quote := q.quotes[rand.Intn(q.totalQuotes)]
	return domain.Quote{
		Author: quote[0],
		Quote:  quote[1],
	}
}
