package profile_attach

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Profile_attachHandler struct {
	db *gorm.DB
}

func NewProfile_attachHandler(db *gorm.DB) *Profile_attachHandler {
	return &Profile_attachHandler{db: db}
}
func (u *Profile_attachHandler) ListProfile_attach(c *gin.Context) {
	var profile_attachs []models.Profile_attach

	r := u.db.Table("profile_attach").Find(&profile_attachs)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile_attachs)
}
func (u *Profile_attachHandler) GetProfile_attachHandler(c *gin.Context) {
	var profile_attach models.Profile_attach
	id := c.Param("id")
	r := u.db.Table("profile_attach").Where("id = ?", id).First(&profile_attach)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile_attach not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile_attach)
}

func (u *Profile_attachHandler) CreateProfile_attachHandler(c *gin.Context) {
	var profile_attach models.Profile_attach

	if err := c.ShouldBindJSON(&profile_attach); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&profile_attach)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *Profile_attachHandler) DeleteProfile_attachHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง profile_attachs
		if err := u.db.Exec("TRUNCATE TABLE profile_attach CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER SEQUENCE profile_attach RESTART WITH 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All profile_attach data have been deleted."})
		return
	}

	// ลบข้อมูล profile_attach ตาม id ที่ระบุ
	r := u.db.Delete(&models.Profile_attach{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("profile_attach with id %s has been deleted.", id)})
}

func (u *Profile_attachHandler) UpdateProfile_attachHandler(c *gin.Context) {
	var profile_attach models.Profile_attach
	id := c.Param("id")

	//ตรวจสอบว่ามี profile_attach นี้อยู่หรือไม่
	r := u.db.Table("profile_attach").Where("id = ?", id).First(&profile_attach)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile_attach not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&profile_attach); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล profile_attach ด้วย ID ที่กำหนด
	r = u.db.Table("profile_attach").Where("id = ?", id).Updates(&profile_attach)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
