package assessment_report

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
	db *gorm.DB
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{db: db}
}
func (u *ReportHandler) ListReport(c *gin.Context) {
	var assessmentreport []models.Report

	r := u.db.Table("assessment_report").Find(&assessmentreport)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessmentreport)
}
func (u *ReportHandler) GetReportHandler(c *gin.Context) {
	var assessmentreport models.Report
	id := c.Param("id")
	r := u.db.Table("assessment_report").Where("id = ?", id).First(&assessmentreport)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment_report not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessmentreport)
}

func (u *ReportHandler) CreateReportHandler(c *gin.Context) {
	var assessmentreport models.Report

	if err := c.ShouldBindJSON(&assessmentreport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&assessmentreport)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *ReportHandler) DeleteReportHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง assessment_report
		if err := u.db.Exec("TRUNCATE TABLE assessment_report CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER TABLE profile ALTER COLUMN id SET DEFAULT nextval('assessment_report_id_seq'::regclass)").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All assessment_report data have been deleted."})
		return
	}

	// ลบข้อมูล assessment_report ตาม id ที่ระบุ
	r := u.db.Delete(&models.Report{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("assessment_report with id %s has been deleted.", id)})
}

func (u *ReportHandler) UpdateReportHandler(c *gin.Context) {
	var assessmentreport models.Report
	id := c.Param("id")

	//ตรวจสอบว่ามี assessment_report นี้อยู่หรือไม่
	r := u.db.Table("assessment_report").Where("id = ?", id).First(&assessmentreport)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment_report not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&assessmentreport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล assessment_report ด้วย ID ที่กำหนด
	r = u.db.Table("assessment_report").Where("id = ?", id).Updates(&assessmentreport)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
