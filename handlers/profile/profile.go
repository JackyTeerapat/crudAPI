package profile

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	db *gorm.DB
}

func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}
func (u *ProfileHandler) ListProfile(c *gin.Context) {
	var profiles []models.Profile

	r := u.db.Table("profile").Find(&profiles)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}
func (u *ProfileHandler) GetProfileHandler(c *gin.Context) {
	var profile models.Profile
	id := c.Param("id")
	r := u.db.Table("profile").Where("id = ?", id).First(&profile)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (u *ProfileHandler) CreateProfileHandler(c *gin.Context) {
	var profile models.Profile

	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&profile)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *ProfileHandler) DeleteProfileHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง profiles
		if err := u.db.Exec("TRUNCATE TABLE profile RESTART IDENTITY CASCADE").Error; err != nil {		
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER SEQUENCE profile RESTART WITH 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All profile data have been deleted."})
		return
	}

	// ลบข้อมูล profile ตาม id ที่ระบุ
	r := u.db.Delete(&models.Profile{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("profile with id %s has been deleted.", id)})
}

func (u *ProfileHandler) UpdateProfileHandler(c *gin.Context) {
	var profile models.Profile
	id := c.Param("id")

	//ตรวจสอบว่ามี profile นี้อยู่หรือไม่
	r := u.db.Table("profile").Where("id = ?", id).First(&profile)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล profile ด้วย ID ที่กำหนด
	r = u.db.Table("profile").Where("id = ?", id).Updates(&profile)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}

