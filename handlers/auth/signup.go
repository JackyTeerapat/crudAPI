package auth

import (
	"CRUD-API/api"
	"CRUD-API/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
