package service

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(n int) (string, error) {
	code := make([]byte, n)
	newInt := big.NewInt(int64(len(alphabet)))
	for i := range code {
		idx, err := rand.Int(rand.Reader, newInt)
		if err != nil {
			return "", err
		}
		code[i] = alphabet[idx.Int64()]
	}
	return string(code), nil
}
