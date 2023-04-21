package assessment

import (
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

	r := u.db.Table("assessment").
		Preload("Project").
		Preload("Progress").
		Preload("Report").
		Preload("Article").
		Find(&assessment)
	if err := r.Error; err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := api.ResponseApiWithDescription(http.StatusOK, assessment, "SUCCESS", nil)
	c.JSON(http.StatusOK, res)
	return
}

func (u *AssessmentHandler) GetAssessmentHandler(c *gin.Context) {
	var assessment models.Assessment
	id := c.Param("id")
	r := u.db.Table("assessment").
		Where("id = ?", id).
		Preload("Project").
		Preload("Progress").
		Preload("Report").
		Preload("Article").
		First(&assessment)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}
	if err := r.Error; err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res := api.ResponseApiWithDescription(http.StatusOK, assessment, "SUCCESS", nil)
	c.JSON(http.StatusOK, res)
	return
}

func (h *AssessmentHandler) CreateAssessmentHandler(c *gin.Context) {
	var assessment models.AssessmentRequest

	if err := c.ShouldBindJSON(&assessment); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Create createdBy and updatedBy variables
	createdBy := "Sahatsawat"
	updatedBy := "Sahatsawat"
	FileName := ""
	FileAction := ""

	// The rest of the code remains the same until the INSERT statements for the other tables

	// Save Project data
	if err := h.db.Exec("INSERT INTO assessment_project (project_year, project_title, project_point, project_estimate, project_recommend, period, created_by, updated_by, file_name, file_action) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		assessment.Project.ProjectYear, assessment.Project.ProjectTitle, assessment.Project.ProjectPoint, assessment.Project.ProjectEstimate, assessment.Project.ProjectRecommend, assessment.Project.ProjectPeriod, createdBy, updatedBy, FileName, FileAction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting project data: %v", err.Error())})
		return
	}

	// Save Progress data
	if err := h.db.Exec("INSERT INTO assessment_progress (progress_year, progress_title, progress_estimate, progress_recommend, period, created_by, updated_by, file_name, file_action) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		assessment.Progress.ProgressYear, assessment.Progress.ProgressTitle, assessment.Progress.ProgressEstimate, assessment.Progress.ProgressRecommend, assessment.Progress.ProgressPeriod, createdBy, updatedBy, FileName, FileAction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting progress data: %v", err.Error())})
		return
	}

	// Save Report data
	if err := h.db.Exec("INSERT INTO assessment_report (report_year, report_title, report_estimate, report_recommend, period, created_by, updated_by, file_name, file_action) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		assessment.Report.ReportYear, assessment.Report.ReportTitle, assessment.Report.ReportEstimate, assessment.Report.ReportRecommend, assessment.Report.ReportPeriod, createdBy, updatedBy, FileName, FileAction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting report data: %v", err.Error())})
		return
	}

	// Save Article data
	if err := h.db.Exec("INSERT INTO assessment_article (article_year, article_title, article_estimate, article_recommend, period, created_by, updated_by, file_name, file_action) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		assessment.Article.ArticleYear, assessment.Article.ArticleTitle, assessment.Article.ArticleEstimate, assessment.Article.ArticleRecommend, assessment.Article.ArticlePeriod, createdBy, updatedBy, FileName, FileAction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting article data: %v", err.Error())})
		return
	}

	var project models.AssessmentProject
	h.db.Last(&project)
	projectID := project.Id

	var progress models.AssessmentProgress
	h.db.Last(&progress)
	progressID := progress.Id

	var report models.AssessmentReport
	h.db.Last(&report)
	reportID := report.Id

	var article models.AssessmentArticle
	h.db.Last(&article)
	articleID := article.Id

	//Assessment
	assessmentFileName := ""
	assessmentFileAction := ""
	// Update the INSERT statement for the assessment table
	result := h.db.Exec("INSERT INTO assessment (profile_id, assessment_start, assessment_end, project_id, progress_id, report_id, article_id, created_by, updated_by, assessment_file_name, assessment_file_action) VALUES (?, ?, ?, ?, ?, ?, ? , ? , ?, ?, ?)",
		assessment.ProfileID, assessment.AssessmentStart, assessment.AssessmentEnd, projectID, progressID, reportID, articleID, createdBy, updatedBy, assessmentFileName, assessmentFileAction)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting assessment data into the assessment table: %v", result.Error)})
		return
	}

	var assessmentid models.Assessment
	h.db.Last(&assessmentid)
	assessmentID := assessmentid.Id

	//------------response------------

	c.JSON(http.StatusCreated, gin.H{"Suscess": fmt.Sprintf("assessment ID : %v Created", assessmentID)})

	var assessmentresponse models.AssessmentResponse
	r := h.db.Table("assessment").
		Where("id = ?", assessmentID).
		Preload("Project").
		Preload("Progress").
		Preload("Report").
		Preload("Article").
		First(&assessmentresponse)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessmentresponse not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := api.ResponseApiWithDescription(http.StatusCreated, assessmentresponse, "CREATED SUCCESS", nil)
	c.JSON(http.StatusCreated, res)
	return

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

		res := api.ResponseApi(http.StatusOK, "deleted", nil)
		c.JSON(http.StatusOK, res)
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
	var assessmentRequest models.AssessmentRequests
	id := c.Param("id")
	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&assessmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล assessment ด้วย ID ที่กำหนด
	var assessment models.Assessment
	u.db.Table("assessment").Where("id = ?", id).First(&assessment)
	result, err := u.update(id, assessmentRequest, assessment)
	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, fmt.Errorf("database error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := api.ResponseApi(http.StatusOK, result, nil)
	c.JSON(http.StatusOK, res)
}

func (u *AssessmentHandler) update(id string, assessmentRequest models.AssessmentRequests, assessment models.Assessment) (body models.Assessment, err error) {
	body = models.Assessment{
		Id:                   assessment.Id,
		Assessment_start:     assessmentRequest.Assessment_start,
		Assessment_end:       assessmentRequest.Assessment_end,
		Assessment_file_name: assessmentRequest.Assessment_file_name,
		Project: models.AssessmentProject{
			Id:                assessment.ProjectID,
			Project_year:      assessmentRequest.Project_year,
			Project_title:     assessmentRequest.Project_title,
			Project_estimate:  assessmentRequest.Project_estimate,
			Project_recommend: assessmentRequest.Project_recommend,
			File_name:         assessmentRequest.Project_file,
			Period:            assessmentRequest.Project_period,
			Updated_by:        "admin",
		},
		Progress: models.AssessmentProgress{
			Id:                 assessment.ProgressID,
			Progress_year:      assessmentRequest.Progress_year,
			Progress_title:     assessmentRequest.Progress_title,
			Progress_estimate:  assessmentRequest.Progress_estimate,
			Progress_recommend: assessmentRequest.Progress_recommend,
			File_name:          assessmentRequest.Progress_file,
			Period:             assessmentRequest.Progress_period,
			Updated_by:         "admin",
		},
		Report: models.AssessmentReport{
			Id:               assessment.ReportID,
			Report_year:      assessmentRequest.Report_year,
			Report_title:     assessmentRequest.Report_title,
			Report_estimate:  assessmentRequest.Report_estimate,
			Report_recommend: assessmentRequest.Report_recommend,
			File_name:        assessmentRequest.Report_file,
			Period:           assessmentRequest.Report_period,
			Updated_by:       "admin",
		},
		Article: models.AssessmentArticle{
			Id:                assessment.ArticleID,
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
	result := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Updates(&body)
	if err := result.Error; err != nil {
		return body, err
	}

	return body, err

}
func (u *AssessmentHandler) create(assessmentRequest models.AssessmentRequests) (body models.Assessment, err error) {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var profile models.Profile

	if err := tx.Find(&profile, 1).Error; err != nil {
		tx.Rollback()
		return body, err
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
		tx.Rollback()
		return body, err
	}

	progress := models.AssessmentProgress{
		Progress_year:      assessmentRequest.Project_year,
		Progress_title:     assessmentRequest.Project_title,
		Progress_estimate:  assessmentRequest.Project_estimate,
		Progress_recommend: assessmentRequest.Project_recommend,
		File_name:          assessmentRequest.Project_file,
		Period:             assessmentRequest.Project_period,
	}

	if err := tx.Create(&progress).Error; err != nil {
		tx.Rollback()
		return body, err
	}
	report := models.AssessmentReport{
		Report_year:      assessmentRequest.Project_year,
		Report_title:     assessmentRequest.Project_title,
		Report_estimate:  assessmentRequest.Project_estimate,
		Report_recommend: assessmentRequest.Project_recommend,
		File_name:        assessmentRequest.Project_file,
		Period:           assessmentRequest.Project_period,
	}

	if err := tx.Create(&report).Error; err != nil {
		tx.Rollback()
		return body, err
	}
	article := models.AssessmentArticle{
		Article_year:      assessmentRequest.Project_year,
		Article_title:     assessmentRequest.Project_title,
		Article_estimate:  assessmentRequest.Project_estimate,
		Article_recommend: assessmentRequest.Project_recommend,
		File_name:         assessmentRequest.Project_file,
		Period:            assessmentRequest.Project_period,
	}

	if err := tx.Create(&article).Error; err != nil {
		tx.Rollback()
		return body, err
	}
	body = models.Assessment{
		Assessment_start:     assessmentRequest.Assessment_start,
		Assessment_end:       assessmentRequest.Assessment_end,
		Assessment_file_name: assessmentRequest.Assessment_file_name,
		ProjectID:            project.Id,
		ProgressID:           progress.Id,
		ProfileID:            profile.ID,
		ReportID:             report.Id,
		ArticleID:            article.Id,
		Created_by:           "admin",
		Updated_by:           "admin",
	}

	r := tx.Debug().Create(&body)
	if err := r.Error; err != nil {
		tx.Rollback()
		return body, err
	}
	return body, tx.Commit().Error
}
