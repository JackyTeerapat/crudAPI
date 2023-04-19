package auth

import (
	"CRUD-API/api"
	"CRUD-API/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (u *AuthHandler) ChangePassword(c *gin.Context) {

	var body struct {
		User_id      uint
		Old_password string
		New_password string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("invalid body"))
		c.JSON(http.StatusBadRequest, res)
		return
	}
	var user models.User
	result := u.db.First(&user, "ID = ?", body.User_id)
	if result.Error != nil {
		res := api.ResponseApiWithDescription(http.StatusBadRequest, nil, "FAILED", result.Error)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Old_password))
	if err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("invalid password or username"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	//Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.New_password), 10)
	if err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("failed to hash password"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Update the "password" column in the database
	u.db.Model(&user).Update("password", string(hash))

	res := api.ResponseApi(http.StatusOK, user, nil)
	c.JSON(http.StatusOK, res)

}
