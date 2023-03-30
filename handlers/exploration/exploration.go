package exploration

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExplorationHandler struct {
	db *gorm.DB
}

func NewExplorationHandler(db *gorm.DB) *ExplorationHandler {
	return &ExplorationHandler{db: db}
}
func (u *ExplorationHandler) ListExploration(c *gin.Context) {
	var explorations []models.Exploration

	r := u.db.Table("exploration").Find(&explorations)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, explorations)
}
func (u *ExplorationHandler) GetExplorationHandler(c *gin.Context) {
	var exploration models.Exploration
	id := c.Param("id")
	r := u.db.Table("exploration").Where("id = ?", id).First(&exploration)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "exploration not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exploration)
}

func (u *ExplorationHandler) CreateExplorationHandler(c *gin.Context) {
	var exploration models.Exploration

	if err := c.ShouldBindJSON(&exploration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&exploration)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *ExplorationHandler) DeleteExplorationHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง explorations
		if err := u.db.Exec("TRUNCATE TABLE exploration CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("SELECT setval('exploration_id_seq', 1, false)").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All exploration data have been deleted."})
		return
	}

	// ลบข้อมูล exploration ตาม id ที่ระบุ
	r := u.db.Delete(&models.Exploration{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("exploration with id %s has been deleted.", id)})
}

func (u *ExplorationHandler) UpdateExplorationHandler(c *gin.Context) {
	var exploration models.Exploration
	id := c.Param("id")

	//ตรวจสอบว่ามี exploration นี้อยู่หรือไม่
	r := u.db.Table("exploration").Where("id = ?", id).First(&exploration)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "exploration not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&exploration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล exploration ด้วย ID ที่กำหนด
	r = u.db.Table("exploration").Where("id = ?", id).Updates(&exploration)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
