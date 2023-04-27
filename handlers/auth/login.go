package auth

import (
	"CRUD-API/api"
	"CRUD-API/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func (u *AuthHandler) Login(c *gin.Context) {

	var body struct {
		Username string
		Password string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("invalid body"))
		c.JSON(http.StatusBadRequest, res)
		return
	}
	v := validator.New()
	regexErr := v.Var(body.Password, "required,min=4,max=8,hexadecimal")
	if regexErr != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("password must be 4-8 character and a-z, A-Z, 0-9"))
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

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("invalid password or username"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

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
