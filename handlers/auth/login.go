package auth

import (
	"CRUD-API/api"
	"CRUD-API/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
