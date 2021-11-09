package entity

import (
	"encoding/base64"
	"math/rand"
	"time"
)

func makePseudoUUID() string {
	// Generate 9 random bytes to obtain 12 base64 characters
	b := make([]byte, 9) //nolint:gomnd

	rand.Seed(time.Now().UnixNano())
	rand.Read(b) //nolint:gosec

	return base64.StdEncoding.EncodeToString(b)
}
