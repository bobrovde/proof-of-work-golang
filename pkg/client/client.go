package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"time"

	"pow/domain"

	jsoniter "github.com/json-iterator/go"
)

const (
	endpointChallenge = "/v1/challenge"
	endpointQuote     = "/v1/quote"

	powHeaderKey = "x-pow-solve"
)

func NewQuoteClient(host string) QuoteClient {
	return QuoteClient{
		host: host,
		httpClient: http.Client{
			Timeout: time.Second * 5,
		},
	}
}

type QuoteClient struct {
	host       string
	httpClient http.Client
}

func (c QuoteClient) GetQuote() (*domain.Quote, error) {
	challenge, err := c.getChallenge()
	if err != nil {
		return nil, err
	}

	nonce := c.solveChallenge(challenge)

	powHeader := getPOWHeader(nonce, challenge.ZeroBytes, challenge.Timestamp, challenge.MAC, challenge.ChallengeData)

	return c.getQuote(powHeader)
}

func (c QuoteClient) getChallenge() (*domain.Challenge, error) {
	request, err := http.NewRequest(http.MethodGet, c.host+endpointChallenge, nil)
	if err != nil {
		return nil, err
	}

	b, err := c.do(request)
	if err != nil {
		return nil, err
	}

	challenge := &domain.Challenge{}

	err = jsoniter.Unmarshal(b, challenge)

	return challenge, err
}

func (c QuoteClient) solveChallenge(challenge *domain.Challenge) int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-challenge.ZeroBytes*8))

	var compareHash big.Int
	var hash [32]byte

	var nonce int64
	data := make([]byte, 0, len(challenge.ChallengeData))

	for nonce < math.MaxInt64 {
		data = append(data, challenge.ChallengeData...)
		data = append(data, intToBytes(nonce)...)

		hash = sha256.Sum256(data)
		compareHash.SetBytes(hash[:])

		if compareHash.Cmp(target) == -1 {
			break
		}
		nonce++
		data = data[:0]
	}
	return int(nonce)
}

func (c QuoteClient) getQuote(powHeader string) (*domain.Quote, error) {
	request, err := http.NewRequest(http.MethodGet, c.host+endpointQuote, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set(powHeaderKey, powHeader)

	b, err := c.do(request)
	if err != nil {
		return nil, err
	}

	quote := &domain.Quote{}

	err = jsoniter.Unmarshal(b, quote)

	return quote, err
}

func (c QuoteClient) do(request *http.Request) ([]byte, error) {
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(string(b))
	}

	return b, nil
}

func getPOWHeader(nonce int, zeroBytesCount int, timestamp int64, mac, challengeData []byte) string {
	return fmt.Sprintf("%d:%d:%d:%s:%s",
		nonce,
		zeroBytesCount,
		timestamp,
		base64.StdEncoding.EncodeToString(mac),
		base64.StdEncoding.EncodeToString(challengeData),
	)
}

func intToBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	_ = binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}
