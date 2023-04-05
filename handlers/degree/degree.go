package degree

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DegreeHandler struct {
	db *gorm.DB
}

func NewDegreeHandler(db *gorm.DB) *DegreeHandler {
	return &DegreeHandler{db: db}
}
func (u *DegreeHandler) ListDegree(c *gin.Context) {
	var degrees []models.Degree

	r := u.db.Table("degree").Find(&degrees)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, degrees)
}
func (u *DegreeHandler) GetDegreeHandler(c *gin.Context) {
	var degree models.Degree
	id := c.Param("id")
	r := u.db.Table("degree").Where("id = ?", id).First(&degree)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "degree not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, degree)
}

func (u *DegreeHandler) CreateDegreeHandler(c *gin.Context) {
	var degree models.Degree

	if err := c.ShouldBindJSON(&degree); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	r := u.db.Create(&degree)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *DegreeHandler) DeleteDegreeHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง degrees
		if err := u.db.Exec("TRUNCATE TABLE degree CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("SELECT setval('degree_id_seq', 1, false)").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All degree data have been deleted."})
		return
	}

	// ลบข้อมูล degree ตาม id ที่ระบุ
	r := u.db.Delete(&models.Degree{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("degree with id %s has been deleted.", id)})
}

func (u *DegreeHandler) UpdateDegreeHandler(c *gin.Context) {
	var degree models.Degree
	id := c.Param("id")

	//ตรวจสอบว่ามี degree นี้อยู่หรือไม่
	r := u.db.Table("degree").Where("id = ?", id).First(&degree)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "degree not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&degree); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล degree ด้วย ID ที่กำหนด
	r = u.db.Table("degree").Where("id = ?", id).Updates(&degree)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
