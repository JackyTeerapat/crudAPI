package assessment

import (
	"fmt"
	"net/http"

	"CRUD-API/api"
	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProgressHandler struct {
	db *gorm.DB
}

func NewProgressHandler(db *gorm.DB) *ProgressHandler {
	return &ProgressHandler{db: db}
}
func (u *ProgressHandler) ListProgress(c *gin.Context) {
	var assessmentprogresses []models.Progress

	r := u.db.Table("assessment_progress").Find(&assessmentprogresses)
	if err := r.Error; err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, assessmentprogresses)
}
func (u *ProgressHandler) GetProgressHandler(c *gin.Context) {
	var assessmentprogress models.Progress
	id := c.Param("id")
	r := u.db.Table("assessment_progress").Where("id = ?", id).First(&assessmentprogress)
	if r.RowsAffected == 0 {
		res := api.ResponseApi(http.StatusNotFound, nil, fmt.Errorf("assessment_progress not found"))
		c.JSON(http.StatusNotFound, res)
		return
	}
	if err := r.Error; err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, assessmentprogress)
}

func (u *ProgressHandler) CreateProgressHandler(c *gin.Context) {
	var assessmentprogress models.Progress

	if err := c.ShouldBindJSON(&assessmentprogress); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("invalid body"))
		c.JSON(http.StatusBadRequest, res)
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&assessmentprogress)
	if err := r.Error; err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *ProgressHandler) DeleteProgressHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง assessment_progress
		if err := u.db.Exec("TRUNCATE TABLE assessment_progress CASCADE").Error; err != nil {
			res := api.ResponseApi(http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, res)
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER TABLE profile ALTER COLUMN id SET DEFAULT nextval('assessment_progress_id_seq'::regclass)").Error; err != nil {
			res := api.ResponseApi(http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, res)
			return
		}
		cp := "deleted all"
		res := api.ResponseApi(http.StatusOK, cp, nil)
		c.JSON(http.StatusOK, res)
		return
	}

	// ลบข้อมูล assessment_progress ตาม id ที่ระบุ
	r := u.db.Delete(&models.Progress{}, id)
	if err := r.Error; err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	cp := "deleted"
	res := api.ResponseApi(http.StatusOK, cp, nil)
	c.JSON(http.StatusOK, res)
}

func (u *ProgressHandler) UpdateProgressHandler(c *gin.Context) {
	var assessmentprogress models.Progress
	id := c.Param("id")

	//ตรวจสอบว่ามี assessment_progress นี้อยู่หรือไม่
	r := u.db.Table("assessment_progress").Where("id = ?", id).First(&assessmentprogress)
	if r.RowsAffected == 0 {
		res := api.ResponseApi(http.StatusNotFound, nil, fmt.Errorf("assessment_progress not found"))
		c.JSON(http.StatusNotFound, res)
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&assessmentprogress); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("invalid body"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	//อัปเดตข้อมูล assessment_progress ด้วย ID ที่กำหนด
	r = u.db.Table("assessment_progress").Where("id = ?", id).Updates(&assessmentprogress)
	if err := r.Error; err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := api.ResponseApi(http.StatusOK, assessmentprogress, nil)
	c.JSON(http.StatusOK, res)
}
