package experience

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExperienceHandler struct {
	db *gorm.DB
}

func NewExperienceHandler(db *gorm.DB) *ExperienceHandler {
	return &ExperienceHandler{db: db}
}
func (u *ExperienceHandler) ListExperience(c *gin.Context) {
	var experiences []models.Experience

	r := u.db.Table("experience").Find(&experiences)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, experiences)
}
func (u *ExperienceHandler) GetExperienceHandler(c *gin.Context) {
	var experience models.Experience
	id := c.Param("id")
	r := u.db.Table("experience").Where("id = ?", id).First(&experience)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "experience not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, experience)
}

func (u *ExperienceHandler) CreateExperienceHandler(c *gin.Context) {
	var experience models.Experience

	if err := c.ShouldBindJSON(&experience); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&experience)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *ExperienceHandler) DeleteExperienceHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง experiences
		if err := u.db.Exec("TRUNCATE TABLE experience CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("SELECT setval('experience_id_seq', 1, false)").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All experience data have been deleted."})
		return
	}

	// ลบข้อมูล experience ตาม id ที่ระบุ
	r := u.db.Delete(&models.Experience{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("experience with id %s has been deleted.", id)})
}

func (u *ExperienceHandler) UpdateExperienceHandler(c *gin.Context) {
	var experience models.Experience
	id := c.Param("id")

	//ตรวจสอบว่ามี experience นี้อยู่หรือไม่
	r := u.db.Table("experience").Where("id = ?", id).First(&experience)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "experience not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&experience); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล experience ด้วย ID ที่กำหนด
	r = u.db.Table("experience").Where("id = ?", id).Updates(&experience)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
