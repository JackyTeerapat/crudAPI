package user

import (
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}
func (u *UserHandler) ListUser(c *gin.Context) {
	var users []models.User

	r := u.db.Table("users").Find(&users)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
func (u *UserHandler) GetUserHandler(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	r := u.db.Table("users").Where("id = ?", id).First(&user)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *UserHandler) CreateUserHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&user)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *UserHandler) DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")

	//รับมาแล้วสร้างเป็น ข้อมูล ลง Tables
	r := u.db.Delete(&models.User{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Status": "Success"})
}
