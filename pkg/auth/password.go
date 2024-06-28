package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(bytes []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(bytes, 10)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func VerifyPassword(hashed []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashed, password)
	return err == nil
}
