package auth

import (
	"CRUD-API/api"
	"CRUD-API/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func (u *AuthHandler) SignUp(c *gin.Context) {
	var body models.Register
	body.Password = "88888888"
	body.Role = "User"

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
	//check username is exist
	var userExist models.User
	u.db.Where("username = ?", body.Username).First(&userExist)
	if userExist.ID > 0 {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("username has already been exist"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	//Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("failed to hash password"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Create user
	user := models.User{Username: body.Username, Password: string(hash), Role: strings.ToUpper(body.Role)}
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
		Password: body.Password,
		Role:     user.Role,
	}

	query := `
	SELECT pg_terminate_backend(pg_stat_activity.pid)
	FROM pg_stat_activity
	WHERE pg_stat_activity.usename = 'navjsbdt'
	AND pg_stat_activity.state = 'idle';
`
	err2 := u.db.Exec(query)
	if err2 != nil {
		fmt.Printf("Failed to close idle connections: %v\n", err2)
	}
	res := api.ResponseApi(http.StatusCreated, response, nil)
	c.JSON(http.StatusCreated, res)
}
