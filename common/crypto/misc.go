package crypto

import (
	"crypto/rand"
	"math/big"

	"github.com/pborman/uuid"
)

// RandNumber returns random number.
func RandNumber(length *big.Int) (*big.Int, error) {
	return rand.Int(rand.Reader, length)
}

// NewUUID returns new uuid.
func NewUUID() string {
	return uuid.NewUUID().String()
}
