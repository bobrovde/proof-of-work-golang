package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"pow/app"
	"pow/pkg/challenge"
	"pow/pkg/quote"
)

func main() {
	quoter, err := quote.NewCSVQuoter()
	if err != nil {
		log.Fatal(err)
	}

	countOfZeroBytes := os.Getenv("COUNT_OF_ZERO_BYTES")
	countOfZeroBytesInt, err := strconv.Atoi(countOfZeroBytes)
	if err != nil {
		log.Fatal(err)
	}

	secretKey := os.Getenv("MAC_SECRET_KEY")

	challengeTTL := os.Getenv("CHALLENGE_TTL")
	chalengeTTLDuration, err := time.ParseDuration(challengeTTL)
	if err != nil {
		log.Fatal(err)
	}

	challenger := challenge.NewChallenger(countOfZeroBytesInt, secretKey, chalengeTTLDuration)

	a := app.NewApp(quoter, challenger)

	httpPort := os.Getenv("HTTP_PORT")
	address := ":" + httpPort

	a.Run(address)
}
