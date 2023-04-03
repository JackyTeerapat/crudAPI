package auth

import (
	"CRUD-API/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	db *gorm.DB
}

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (u *AuthHandler) SignUp(c *gin.Context) {
	var body struct {
		Username string
		Password string
		Role     string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to hash password"})
		return
	}

	// Create user
	user := models.Register{Username: body.Username, Password: string(hash), Role: body.Role}
	res := u.db.Create(&user)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}

	//Response
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (u *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Username string
		Password string
		Role     string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	var user models.Register
	u.db.First(&user, "username = ?", body.Username)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := GenerateToken()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":      http.StatusOK,
			"description": "SUCCESS",
			"data": gin.H{
				"user_id":  user.ID,
				"username": user.Username,
				"role":     user.Role,
				"token":    token,
			},
		},
	)
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
		Issuer:    "test",
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

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secretKey, nil
	})

	return
}
