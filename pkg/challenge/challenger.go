package challenge

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/big"
	"math/rand"
	"time"

	"pow/domain"
)

type Challenger struct {
	zeroBytesCount int
	target         *big.Int
	secretKey      []byte
	challengeTTL   time.Duration
}

func NewChallenger(zeroBytesCount int, secretKey string, challengeTTL time.Duration) *Challenger {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-zeroBytesCount*8))

	return &Challenger{
		zeroBytesCount: zeroBytesCount,
		target:         target,
		secretKey:      []byte(secretKey),
		challengeTTL:   challengeTTL,
	}
}

func (c Challenger) GetChallenge() (*domain.Challenge, error) {
	challengeData := make([]byte, 16)
	_, err := rand.Read(challengeData)
	if err != nil {
		return nil, err
	}

	nowUnix := time.Now().UnixNano()

	mac := c.generateMAC(challengeData, nowUnix)

	return &domain.Challenge{
		ChallengeData: challengeData,
		ZeroBytes:     c.zeroBytesCount,
		MAC:           mac,
		Timestamp:     nowUnix,
	}, nil
}

// If need guranty that challengeData used once we can add map to store usage of challenge data
func (c *Challenger) ValidateChallengeSolve(nonce int, unixNano int64, challengeData, mac []byte) error {
	expectedMac := c.generateMAC(challengeData, unixNano)
	if !hmac.Equal(expectedMac, mac) {
		return errors.New("wrong challenge solution")
	}

	if time.Now().Sub(time.Unix(0, unixNano)) > c.challengeTTL {
		return errors.New("challenge exceeded")
	}

	data := append(challengeData, intToBytes(int64(nonce))...)
	challengeHash := sha256.Sum256(data)

	var challengeHashInt big.Int
	challengeHashInt.SetBytes(challengeHash[:])

	if challengeHashInt.Cmp(c.target) != -1 {
		return errors.New("wrong challenge solution")
	}

	return nil
}

func (c *Challenger) generateMAC(challengeData []byte, timestamp int64) []byte {
	mac := hmac.New(sha256.New, c.secretKey)
	mac.Write(challengeData)
	mac.Write(intToBytes(timestamp))
	return mac.Sum(nil)
}

func intToBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	_ = binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}
