package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var hmacSampleSecret = []byte("3HqC0s7224Ai0XTnPBhURakeGjPaPNcUWOtrz9N+hpQ=")

func GenerateToken() (tokenString string, err error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
		Issuer:    "test",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString(hmacSampleSecret)

	return
}

func ValidateToken(tokenString string) (token *jwt.Token, err error) {

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	return
}
