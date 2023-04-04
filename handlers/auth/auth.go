package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	db *gorm.DB
}

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func GenerateToken(id uint, username string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":      jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		"username": username,
		"id":       strconv.Itoa(int(id)),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString(secretKey)

	if err != nil {
		return "fail to create token", err
	}

	return
}

func ValidateToken(tokenString string) (token *jwt.Token, err error) {

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	return
}
