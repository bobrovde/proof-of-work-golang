package app

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (a *App) quoteV1(c echo.Context) error {
	challengeHeader := c.Request().Header.Get("x-pow-solve")
	if challengeHeader == "" {
		return errors.New("empty header x-pow-solve")
	}

	fmt.Println("got x-pow-solve header", challengeHeader)

	challengeSolveData, err := extractChallengeSolveData(challengeHeader)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	solveErr := a.challenger.ValidateChallengeSolve(
		challengeSolveData.Nonce,
		challengeSolveData.Timestamp,
		challengeSolveData.ChallengeData,
		challengeSolveData.MAC,
	)
	if solveErr != nil {
		return c.String(http.StatusForbidden, solveErr.Error())
	}

	quote := a.quoter.GetRandomQuote()

	return c.JSON(http.StatusOK, quote)
}

type challengeSolveData struct {
	Nonce          int
	ZeroBytesCount int
	ChallengeData  []byte
	Timestamp      int64
	MAC            []byte
}

func extractChallengeSolveData(challengeSolveHeader string) (data challengeSolveData, err error) {
	//nonce:zeroBytesCount:date:mac:challenge
	challengeData := strings.Split(challengeSolveHeader, ":")
	if len(challengeData) < 4 {
		err = errors.New("invalid header format, want nonce:target:data")
		return
	}

	data.Nonce, err = strconv.Atoi(challengeData[0])
	if err != nil {
		err = fmt.Errorf("invalid format of nonce, want int:%s", err)
		return
	}

	data.ZeroBytesCount, err = strconv.Atoi(challengeData[1])
	if err != nil {
		err = fmt.Errorf("invalid format of zero bytes count, want int:%s", err)
		return
	}

	data.Timestamp, err = strconv.ParseInt(challengeData[2], 10, 64)
	if err != nil {
		err = fmt.Errorf("failed to decode timestamp:%s", err)
		return
	}
	data.MAC, err = base64.StdEncoding.DecodeString(challengeData[3])
	if err != nil {
		err = fmt.Errorf("failed to decode mac:%s", err)
		return
	}
	data.ChallengeData, err = base64.StdEncoding.DecodeString(challengeData[4])
	if err != nil {
		err = fmt.Errorf("failed to decode challenge data:%s", err)
		return
	}
	return
}
