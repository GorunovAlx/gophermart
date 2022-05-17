package v1

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Generate a token.
func GenerateUserIDToken() (string, error) {
	id, err := generateRandom(4)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(id)
	signedID := hex.EncodeToString(append(id, h.Sum(nil)...))

	return signedID, nil
}

// Authenticate user token id.
func AuthUserIDToken(userIDToken string) (bool, error) {
	data, err := hex.DecodeString(userIDToken)
	if err != nil {
		return false, err
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data[:4])
	sign := h.Sum(nil)

	if !hmac.Equal(sign, data[4:]) {
		return false, nil
	}

	return true, nil
}
