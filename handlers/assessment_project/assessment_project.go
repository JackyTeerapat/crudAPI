package assessment_project

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssessmentProjectHandler struct {
	db *gorm.DB
}

func NewAssessmentProjectHandler(db *gorm.DB) *AssessmentProjectHandler {
	return &AssessmentProjectHandler{db: db}
}
func (u *AssessmentProjectHandler) ListAssessmentProjects(c *gin.Context) {
	var assessmentProject []models.Assessment_project

	r := u.db.Table("assessment_project").Find(&assessmentProject)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessmentProject)
}
func (u *AssessmentProjectHandler) GetAssessmentProjectHandler(c *gin.Context) {
	var assessmentProject models.Assessment_project
	id := c.Param("id")
	r := u.db.Table("assessment_project").Where("id = ?", id).First(&assessmentProject)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assessment project not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessmentProject)
}

func (u *AssessmentProjectHandler) CreateAssessmentProjectHandler(c *gin.Context) {
	var assessmentProject models.Assessment_project

	if err := c.ShouldBindJSON(&assessmentProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&assessmentProject)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *AssessmentProjectHandler) DeleteAssessmentProjectHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง assessment project
		if err := u.db.Exec("TRUNCATE TABLE assessment_project CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER SEQUENCE assessment_project RESTART WITH 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All assessment project data have been deleted."})
		return
	}

	// ลบข้อมูลตาม id ที่ระบุ
	r := u.db.Delete(&models.Assessment_project{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("assessment project with id %s has been deleted.", id)})
}

func (u *AssessmentProjectHandler) UpdateAssessmentProjectHandler(c *gin.Context) {
	var assessmentProject models.Degree
	id := c.Param("id")

	//ตรวจสอบว่ามี degree นี้อยู่หรือไม่
	r := u.db.Table("assessment_project").Where("id = ?", id).First(&assessmentProject)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment project not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&assessmentProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล degree ด้วย ID ที่กำหนด
	r = u.db.Table("assessment_project").Where("id = ?", id).Updates(&assessmentProject)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
