package domain

import (
	"fmt"
)

type Challenge struct {
	ZeroBytes     int
	ChallengeData []byte
	MAC           []byte
	Timestamp     int64
}

type Quote struct {
	Author string
	Quote  string
}

func (q Quote) String() string {
	authorName := q.Author
	if authorName == "" {
		authorName = "Unknown"
	}
	return fmt.Sprintf("author: %s\nquote: %s\n", authorName, q.Quote)
}
