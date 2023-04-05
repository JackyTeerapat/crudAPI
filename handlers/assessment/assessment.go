package assessment

import (
	"fmt"
	"net/http"

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

	r := u.db.Table("assessment").Find(&assessment)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessment)
}

func (u *AssessmentHandler) GetAssessmentHandler(c *gin.Context) {
	var assessment models.Assessment
	id := c.Param("id")
	r := u.db.Table("assessment").Where("id = ?", id).Preload("Project").Preload("Progress").Preload("Report").Preload("Article").Preload("Profile").First(&assessment)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessment)
}

func (h *AssessmentHandler) CreateAssessmentHandler(c *gin.Context) {
	var assessment models.AssessmentRequest

	if err := c.ShouldBindJSON(&assessment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create createdBy and updatedBy variables
	createdBy := "Champlnwza007"
	updatedBy := "Champlnwza007"
	// The rest of the code remains the same until the INSERT statements for the other tables

	// Save Project data
	if err := h.db.Exec("INSERT INTO assessment_project (project_year, project_title, project_point, project_estimate, project_recommend, period, file_name, file_storage, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		assessment.Project.ProjectYear, assessment.Project.ProjectTitle, assessment.Project.ProjectPoint, assessment.Project.ProjectEstimate, assessment.Project.ProjectRecommend, assessment.Project.ProjectPeriod, assessment.Project.ProjectFile.FileName, assessment.Project.ProjectFile.FileStorage, createdBy, updatedBy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting project data: %v", err.Error())})
		return
	}

	// Save Progress data
	if err := h.db.Exec("INSERT INTO assessment_progress (progress_year, progress_title, progress_estimate, progress_recommend, period, file_name, file_storage, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		assessment.Progress.ProgressYear, assessment.Progress.ProgressTitle, assessment.Progress.ProgressEstimate, assessment.Progress.ProgressRecommend, assessment.Progress.ProgressPeriod, assessment.Progress.ProgressFile.FileName, assessment.Progress.ProgressFile.FileStorage, createdBy, updatedBy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting progress data: %v", err.Error())})
		return
	}

	// Save Report data
	if err := h.db.Exec("INSERT INTO assessment_report (report_year, report_title, report_estimate, report_recommend, period, file_name, file_storage, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		assessment.Report.ReportYear, assessment.Report.ReportTitle, assessment.Report.ReportEstimate, assessment.Report.ReportRecommend, assessment.Report.ReportPeriod, assessment.Report.ReportFile.FileName, assessment.Report.ReportFile.FileStorage, createdBy, updatedBy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting report data: %v", err.Error())})
		return
	}

	// Save Article data
	if err := h.db.Exec("INSERT INTO assessment_article (article_year, article_title, article_estimate, article_recommend, period, file_name, file_storage, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		assessment.Article.ArticleYear, assessment.Article.ArticleTitle, assessment.Article.ArticleEstimate, assessment.Article.ArticleRecommend, assessment.Article.ArticlePeriod, assessment.Article.ArticleFile.FileName, assessment.Article.ArticleFile.FileStorage, createdBy, updatedBy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting article data: %v", err.Error())})
		return
	}

	var project models.Project
	h.db.Last(&project)
	projectID := project.Id

	var progress models.Progress
	h.db.Last(&progress)
	progressID := progress.Id

	var report models.Report
	h.db.Last(&report)
	reportID := report.Id

	var article models.Article
	h.db.Last(&article)
	articleID := article.Id

	var profile models.Profile
	h.db.Last(&profile)
	profileID := profile.ID

	//Assessment

	// Update the INSERT statement for the assessment table
	result := h.db.Exec("INSERT INTO assessment (assessment_start, assessment_end, assessment_file_name, assessment_file_storage, project_id, progress_id, report_id, article_id, profile_id, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ? , ? , ?, ?, ?, ?)",
		assessment.AssessmentStart, assessment.AssessmentEnd, assessment.AssessmentFile.FileName, assessment.AssessmentFile.FileStorage, projectID, progressID, reportID, articleID, profileID, createdBy, updatedBy)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting assessment data into the assessment table: %v", result.Error)})
		return
	}

	var assessmentid models.Assessment
	h.db.Last(&assessmentid)
	assessmentID := assessmentid.Id

	c.JSON(http.StatusCreated, gin.H{"Suscess": fmt.Sprintf("assessment ID : %v Created", assessmentID)})
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
	var assessment models.Assessment
	id := c.Param("id")

	//ตรวจสอบว่ามี assessment นี้อยู่หรือไม่
	r := u.db.Table("assessment").Where("id = ?", id).First(&assessment)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&assessment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล assessment ด้วย ID ที่กำหนด
	r = u.db.Table("assessment").Where("id = ?", id).Updates(&assessment)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
