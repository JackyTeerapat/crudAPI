package assessment

import (
	"errors"
	"fmt"
	"net/http"

	"CRUD-API/api"
	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssessmentHandler struct {
	db *gorm.DB
}

func NewAssessmentHandler(db *gorm.DB) *AssessmentHandler {
	return &AssessmentHandler{db: db}
}
func (u *AssessmentHandler) ListAssessment(c *gin.Context) {
	var assessment []models.Assessment

	r := u.db.Table("assessment").Preload("Project").Preload("Progress").Preload("Report").Preload("Article").Find(&assessment)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessment)
}
func (u *AssessmentHandler) GetAssessmentHandler(c *gin.Context) {
	var assessment models.Assessment
	id := c.Param("id")
	r := u.db.Table("assessment").Preload("Project").Preload("Progress").Preload("Report").Preload("Article").Where("id = ?", id).First(&assessment)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessment)
}
func (u *AssessmentHandler) create(assessmentRequest models.AssessmentRequest) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Debug().Rollback()
		}
	}()
	var profile models.Profile

	if err := tx.Find(&profile, 1).Error; err != nil {
		tx.Rollback()
		return err
	}
	project := models.AssessmentProject{
		Project_year:      assessmentRequest.Project_year,
		Project_title:     assessmentRequest.Project_title,
		Project_point:     assessmentRequest.Project_point,
		Project_estimate:  assessmentRequest.Project_estimate,
		Project_recommend: assessmentRequest.Project_recommend,
		File_name:         assessmentRequest.Project_file,
		Period:            assessmentRequest.Project_period,
	}
	p := tx.Create(&project)
	if err := p.Error; err != nil {
		tx.Debug().Rollback()
		return err
	}

	progress := models.Progress{
		Progress_year:      assessmentRequest.Project_year,
		Progress_title:     assessmentRequest.Project_title,
		Progress_estimate:  assessmentRequest.Project_estimate,
		Progress_recommend: assessmentRequest.Project_recommend,
		File_name:          assessmentRequest.Project_file,
		Period:             assessmentRequest.Project_period,
	}

	if err := tx.Create(&progress).Error; err != nil {
		tx.Debug().Rollback()
		return err
	}
	report := models.Report{
		Report_year:      assessmentRequest.Project_year,
		Report_title:     assessmentRequest.Project_title,
		Report_estimate:  assessmentRequest.Project_estimate,
		Report_recommend: assessmentRequest.Project_recommend,
		File_name:        assessmentRequest.Project_file,
		Period:           assessmentRequest.Project_period,
	}

	if err := tx.Create(&report).Error; err != nil {
		tx.Debug().Rollback()
		return err
	}
	article := models.Article{
		Article_year:      assessmentRequest.Project_year,
		Article_title:     assessmentRequest.Project_title,
		Article_estimate:  assessmentRequest.Project_estimate,
		Article_recommend: assessmentRequest.Project_recommend,
		File_name:         assessmentRequest.Project_file,
		Period:            assessmentRequest.Project_period,
	}

	if err := tx.Create(&article).Error; err != nil {
		tx.Debug().Rollback()
		return err
	}
	data := models.Assessment{
		Assessment_start:     assessmentRequest.Assessment_start,
		Assessment_end:       assessmentRequest.Assessment_end,
		Assessment_file_name: assessmentRequest.Assessment_file_name,
		ProjectID:            project.ID,
		ProgressID:           progress.ID,
		ProfileID:            profile.ID,
		ReportID:             report.ID,
		ArticleID:            article.ID,
		Created_by:           "admin",
		Updated_by:           "admin",
	}

	r := tx.Debug().Create(&data)
	if err := r.Error; err != nil {
		tx.Debug().Rollback()
		return err
	}
	return tx.Debug().Commit().Error
}

func (u *AssessmentHandler) CreateAssessmentHandler(c *gin.Context) {
	var assessment models.AssessmentRequest
	if err := c.ShouldBindJSON(&assessment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	if err := u.create(assessment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	res := api.ResponseApi(200, "Success", nil)
	c.JSON(http.StatusOK, res)
}

func (u *AssessmentHandler) DeleteAssessmentHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง assessment
		if err := u.db.Exec("TRUNCATE TABLE assessment CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER TABLE profile ALTER COLUMN id SET DEFAULT nextval('assessment_id_seq'::regclass)").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All assessment data have been deleted."})
		return
	}

	// ลบข้อมูล assessment ตาม id ที่ระบุ
	r := u.db.Delete(&models.Assessment{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("assessment with id %s has been deleted.", id)})
}

func (u *AssessmentHandler) UpdateAssessmentHandler(c *gin.Context) {
	var assessmentRequest models.AssessmentRequest
	id := c.Param("id")
	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&assessmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var assessment models.Assessment
	// r := tx.Table("assessment").Where("id = ?", id).First(&assessment)
	// if r.RowsAffected == 0 {
	// 	tx.Rollback()
	// 	res := api.ResponseApi(http.StatusNotFound, nil, fmt.Errorf("assessment not found"))
	// 	c.JSON(http.StatusNotFound, res)
	// 	return
	// }
	result := u.db.Table("assessment").Where("id = ?", id).First(&assessment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := models.Assessment{
		Assessment_start:     assessmentRequest.Assessment_start,
		Assessment_end:       assessmentRequest.Assessment_end,
		Assessment_file_name: assessmentRequest.Assessment_file_name,
		Project: models.AssessmentProject{
			ID:                assessment.ProjectID,
			Project_year:      assessmentRequest.Project_year,
			Project_title:     assessmentRequest.Project_title,
			Project_estimate:  assessmentRequest.Project_estimate,
			Project_recommend: assessmentRequest.Project_recommend,
			File_name:         assessmentRequest.Project_file,
			Period:            assessmentRequest.Project_period,
			Updated_by:        "admin",
		},
		Progress: models.Progress{
			ID:                 assessment.ProgressID,
			Progress_year:      assessmentRequest.Progress_year,
			Progress_title:     assessmentRequest.Progress_title,
			Progress_estimate:  assessmentRequest.Progress_estimate,
			Progress_recommend: assessmentRequest.Progress_recommend,
			File_name:          assessmentRequest.Progress_file,
			Period:             assessmentRequest.Progress_period,
			Updated_by:         "admin",
		},
		Report: models.Report{
			ID:               assessment.ReportID,
			Report_year:      assessmentRequest.Report_year,
			Report_title:     assessmentRequest.Report_title,
			Report_estimate:  assessmentRequest.Report_estimate,
			Report_recommend: assessmentRequest.Report_recommend,
			File_name:        assessmentRequest.Report_file,
			Period:           assessmentRequest.Report_period,
			Updated_by:       "admin",
		},
		Article: models.Article{
			ID:                assessment.ArticleID,
			Article_year:      assessmentRequest.Article_year,
			Article_title:     assessmentRequest.Article_title,
			Article_estimate:  assessmentRequest.Article_estimate,
			Article_recommend: assessmentRequest.Article_recommend,
			File_name:         assessmentRequest.Article_file,
			Period:            assessmentRequest.Article_period,
			Updated_by:        "admin",
		},
		Created_by: "admin",
		Updated_by: "admin",
	}
	//อัปเดตข้อมูล assessment ด้วย ID ที่กำหนด
	result = u.db.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Updates(&data)
	if err := result.Error; err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, fmt.Errorf("database error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := api.ResponseApi(http.StatusOK, assessmentRequest, nil)
	c.JSON(http.StatusOK, res)
}
