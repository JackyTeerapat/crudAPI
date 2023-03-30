package user

import (
	"fmt"
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

// GetUser godoc
// @Summary Get a user
// @Description Get a data user from database.
// @Tags User
// @Produce  application/json
// @Param id path int true "User ID"
// @Success 200 {object} models.User{}
// @Router /user/{id} [get]
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

// CreateUser godoc
// @Summary Create a user
// @Description Create a data user to database.
// @Tags User
// @Produce  application/json
// @Param user body models.User true "User"
// @Success 200 {object} models.User{}
// @Router /user [post]
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

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a data user from database.
// @Tags User
// @Produce  application/json
// @Param id path int true "User ID"
// @Success 200 {object} models.User{}
// @Router /user/{id} [delete]
func (u *UserHandler) DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง Users
		if err := u.db.Exec("TRUNCATE TABLE users CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All user data have been deleted."})
		return
	}

	// ลบข้อมูล User ตาม id ที่ระบุ
	r := u.db.Delete(&models.User{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User with id %s has been deleted.", id)})
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a data user to database.
// @Tags User
// @Produce  application/json
// @Param id path int true "User ID"
// @Param user body models.User true "User"
// @Success 200 {object} models.User{}
// @Router /user/{id} [put]
func (u *UserHandler) UpdateUserHandler(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	//ตรวจสอบว่ามี User นี้อยู่หรือไม่
	r := u.db.Table("users").Where("id = ?", id).First(&user)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล User ด้วย ID ที่กำหนด
	r = u.db.Table("users").Where("id = ?", id).Updates(&user)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
