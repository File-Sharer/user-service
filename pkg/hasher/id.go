package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateUserID(login string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(login))
	if err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString[:16], nil
}
