package auth

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	db *gorm.DB
}

var hmacSampleSecret = []byte("3HqC0s7224Ai0XTnPBhURakeGjPaPNcUWOtrz9N+hpQ=")

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (u *AuthHandler) Login(c *gin.Context) {

	token, err := GenerateToken()

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"token": token,
		},
	)

	// var user models.User
	// var userLogin models.User
	// if err := c.ShouldBindJSON(&userLogin); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
	// r := u.db.Table("users").Where("username = ?", userLogin.Username).First(&user)
	// if r.RowsAffected == 0 {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	// 	return
	// }
	// if err := r.Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// if userLogin.Password != user.Password {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
	// 	return
	// }
	// token, err := GenerateToken()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{"token": token})
}

func (u *AuthHandler) GetUser(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"userId": c.Param("id"),
		},
	)
}
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
