package app

import (
	"pow/domain"

	"github.com/labstack/echo/v4"
)

type quoter interface {
	GetRandomQuote() domain.Quote
}

type challenger interface {
	GetChallenge() (*domain.Challenge, error)
	ValidateChallengeSolve(nonce int, unixNano int64, challengeData, mac []byte) error
}

func NewApp(quoter quoter, challenger challenger) *App {
	e := echo.New()

	a := &App{
		e:          e,
		challenger: challenger,
		quoter:     quoter,
	}
	a.registerRoutes()
	return a
}

type App struct {
	e *echo.Echo

	challenger challenger
	quoter     quoter
}

func (a *App) Run(address string) {
	a.e.Logger.Fatal(a.e.Start(address))
}

func (a *App) registerRoutes() {
	a.e.GET("/v1/challenge", a.challengeV1)
	a.e.GET("/v1/quote", a.quoteV1)
}
