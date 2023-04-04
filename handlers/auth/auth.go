package auth

import (
	"CRUD-API/handlers/api"
	"CRUD-API/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	var body models.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("invalid body"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	//check username is exist
	var userExist models.User
	u.db.Where("username = ?", body.Username).First(&userExist)
	if userExist.ID > 0 {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("username is exist"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Hash password
	// hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	// if err != nil {
	// 	res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("failed to hash password"))
	// 	c.JSON(http.StatusBadRequest, res)
	// 	return
	// }
	// Create user
	user := models.User{Username: body.Username, Role: strings.ToUpper(body.Role)}
	data := u.db.Create(&user)
	if data.Error != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, data.Error)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	//Response
	response := models.Register{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	res := api.ResponseApi(http.StatusCreated, response, nil)
	c.JSON(http.StatusCreated, res)
}

func (u *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Username string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("invalid body"))
		c.JSON(http.StatusBadRequest, res)
		return
	}
	var user models.User
	u.db.First(&user, "username = ?", body.Username)
	if user.ID == 0 {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("user not found"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
	// 	return
	// }

	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		res := api.ResponseApi(http.StatusBadGateway, nil, err)
		c.JSON(http.StatusBadGateway, res)
		return
	}
	data := models.LoginRespones{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Token:    token,
	}
	res := api.ResponseApi(http.StatusOK, data, nil)
	c.JSON(http.StatusOK, res)
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
