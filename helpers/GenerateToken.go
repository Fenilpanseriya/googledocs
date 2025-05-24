package helpers

import (
	"time"

	"github.com/fenilpanseriya/docs2.0/models"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(user *models.User, jwtKey string, expireTime time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(expireTime).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {

		return "", err
	}
	return tokenString, nil
}
