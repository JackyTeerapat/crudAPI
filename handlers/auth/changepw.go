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

type CustomResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

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
	v := validator.New()
	regexErr := v.Var(body.New_password, "required,min=4,max=8,hexadecimal")
	if regexErr != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("password must be 4-8 character and a-z, A-Z, 0-9"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if body.Old_password == body.New_password {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("can't use the same password"))
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

	customResponse := CustomResponse{
		UserID:   user.ID,
		Username: user.Username,
		Password: user.Password,
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
	res := api.ResponseApi(http.StatusOK, customResponse, nil)
	c.JSON(http.StatusOK, res)

}
