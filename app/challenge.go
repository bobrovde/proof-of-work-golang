package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *App) challengeV1(c echo.Context) error {
	challenge, err := a.challenger.GetChallenge()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, challenge)
}
