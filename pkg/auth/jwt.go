package auth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	jwt, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func SendToken(c *gin.Context, token string, host string) error {
	c.SetCookie("jwt", token, int(time.Now().Add(time.Hour * 24).Unix()), "/", host, true, true)
	return nil
}

func GetTokenClaims(token string, secretKey []byte) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}
