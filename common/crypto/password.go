package crypto

import (
	"math/big"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"

	"github.com/privatix/dapp-openvpn/common/transform"
)

// HashPassword computes encoded hash of the password.
func HashPassword(password, salt string) (string, error) {
	salted := []byte(password + salt)
	passwordHash, err := bcrypt.GenerateFromPassword(salted,
		bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return transform.FromBytes(passwordHash), nil
}

// RandPass returns random password.
func RandPass(length int) string {
	n, _ := RandNumber(big.NewInt(10))
	pass, _ := password.Generate(length, int(n.Int64()), 0, false, false)
	return pass
}
