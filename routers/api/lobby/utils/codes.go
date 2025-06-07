package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	// The characters you want in your codes:
	charset    = "123456789"
	DefaultLen = 6
)

// GenerateCode returns a random string of length `n` drawn from [a–z1–9].
// It uses crypto/rand so collisions are still extremely unlikely even
// if you generate thousands of them.
func GenerateCode(n int) (string, error) {
	b := make([]byte, n)
	max := big.NewInt(int64(len(charset)))
	for i := range b {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}
