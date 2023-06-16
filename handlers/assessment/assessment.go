package assessment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

type AssessmentResponse struct {
	profile_id int
	Project    []models.AssessmentProjectGet  `json:"Project"`
	Progress   []models.AssessmentProgressGet `json:"Progress"`
	Report     []models.AssessmentReportGet   `json:"Report"`
	Article    []models.AssessmentArticleGet  `json:"Article"`
}

func (u *AssessmentHandler) GetAssessmentHandler(c *gin.Context) {
	id := c.Param("id")
	var project []models.AssessmentProjectGet
	var progress []models.AssessmentProgressGet
	var report []models.AssessmentReportGet
	var article []models.AssessmentArticleGet

	r := u.db.Table("assessment_project").
		Where("profile_id = ?", id).Find(&project)

	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}
	if err := r.Error; err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	r = u.db.Table("assessment_progress").
		Where("profile_id = ?", id).Find(&progress)

	if r.Error != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, r.Error)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	r = u.db.Table("assessment_report").
		Where("profile_id = ?", id).Find(&report)

	if r.Error != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, r.Error)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	r = u.db.Table("assessment_article").
		Where("profile_id = ?", id).Find(&article)

	if r.Error != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, r.Error)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	profileID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	responseData := struct {
		ProfileID int                            `json:"profile_id"`
		Project   []models.AssessmentProjectGet  `json:"Project"`
		Progress  []models.AssessmentProgressGet `json:"Progress"`
		Report    []models.AssessmentReportGet   `json:"Report"`
		Article   []models.AssessmentArticleGet  `json:"Article"`
	}{
		ProfileID: profileID,
		Project:   project,
		Progress:  progress,
		Report:    report,
		Article:   article,
	}

	res := api.ResponseApiWithDescription(http.StatusOK, responseData, "SUCCESS", nil)
	c.JSON(http.StatusOK, res)
	return
}

func (u *AssessmentHandler) CreateAssessmentHandler(c *gin.Context) {
	var assessment models.AssessmentRequests
	if err := c.ShouldBindJSON(&assessment); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	var profile models.Profile
	ckeck := u.db.Table("profile").Where("id = ?", assessment.ProfileID).First(&profile)
	if ckeck.RowsAffected == 0 {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("no data found for this profile"))
		c.JSON(http.StatusBadRequest, res)
		return
	}
	body, err := u.create(assessment)

	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := api.ResponseApi(http.StatusCreated, body, nil)
	c.JSON(http.StatusCreated, res)
}

type UpdateAssessmentRequest struct {
	ProfileID      int    `json:"profile_id"`
	AssessmentType string `json:"assessment_type"`
	ProjectID      int    `json:"project_id,omitempty"`
	ProgressID     int    `json:"progress_id,omitempty"`
	ReportID       int    `json:"report_id,omitempty"`
	ArticleID      int    `json:"article_id,omitempty"`
}

type UpdateAssessmentResponse struct {
	Data         map[string]bool `json:"data"`
	Description  string          `json:"description"`
	ErrorMessage interface{}     `json:"errorMessage"`
	Status       int             `json:"status"`
}

func (u *AssessmentHandler) UpdateAssessmentStatusHandler(c *gin.Context) {
	requestBody := UpdateAssessmentRequest{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		res := UpdateAssessmentResponse{
			Description: "Invalid request body",
			Status:      http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var tableName string
	var statusField string
	var assessmentID int
	switch requestBody.AssessmentType {
	case "project":
		tableName = "assessment_project"
		statusField = "project_status"
		assessmentID = requestBody.ProjectID
	case "progress":
		tableName = "assessment_progress"
		statusField = "progress_status"
		assessmentID = requestBody.ProgressID
	case "report":
		tableName = "assessment_report"
		statusField = "report_status"
		assessmentID = requestBody.ReportID
	case "article":
		tableName = "assessment_article"
		statusField = "article_status"
		assessmentID = requestBody.ArticleID
	default:
		res := UpdateAssessmentResponse{
			Description: "Invalid assessment type",
			Status:      http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if assessmentID == 0 {
		res := UpdateAssessmentResponse{
			Description: "Invalid assessment ID",
			Status:      http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	query := u.db.Table(tableName).
		Where("profile_id = ? AND id = ?", requestBody.ProfileID, assessmentID).
		Update(statusField, false)
	if query.RowsAffected == 0 {
		res := UpdateAssessmentResponse{
			Description: "Assessment not found",
			Status:      http.StatusNotFound,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}
	if query.Error != nil {
		res := UpdateAssessmentResponse{
			Description:  "Database error",
			ErrorMessage: query.Error,
			Status:       http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	responseData := make(map[string]bool)
	responseData[statusField] = false

	res := UpdateAssessmentResponse{
		Data:        responseData,
		Description: "SUCCESS",
		Status:      http.StatusOK,
	}
	c.JSON(http.StatusOK, res)
}

func (u *AssessmentHandler) updateProjectStatus(profileID string, projectID int, status bool) error {
	return u.db.Table("assessment_project").
		Where("profile_id = ? AND id = ?", profileID, projectID).
		Update("project_status", status).Error
}

func (u *AssessmentHandler) updateProgressStatus(profileID string, progressID int, status bool) error {
	return u.db.Table("assessment_progress").
		Where("profile_id = ? AND id = ?", profileID, progressID).
		Update("progress_status", status).Error
}

func (u *AssessmentHandler) updateReportStatus(profileID string, reportID int, status bool) error {
	return u.db.Table("assessment_report").
		Where("profile_id = ? AND id = ?", profileID, reportID).
		Update("report_status", status).Error
}

func (u *AssessmentHandler) updateArticleStatus(profileID string, articleID int, status bool) error {
	return u.db.Table("assessment_article").
		Where("profile_id = ? AND id = ?", profileID, articleID).
		Update("article_status", status).Error
}

func (u *AssessmentHandler) UpdateAssessmentHandler(c *gin.Context) {
	var assessmentRequest models.AssessmentRequests
	id := c.Param("id")
	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&assessmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profileId, err := strconv.Atoi(id)
	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, fmt.Errorf("failed to convert string to int"))
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	//อัปเดตข้อมูล assessment ด้วย ID ที่กำหนด
	var profile models.Profile
	ckeck := u.db.Table("profile").Where("id = ?", assessmentRequest.ProfileID).First(&profile)
	if ckeck.RowsAffected == 0 {
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("no data found for this profile"))
		c.JSON(http.StatusBadRequest, res)
		return
	}
	result, err := u.update(profileId, assessmentRequest, profile)
	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, fmt.Errorf("database error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := api.ResponseApi(http.StatusOK, result, nil)
	c.JSON(http.StatusOK, res)
}

func (u *AssessmentHandler) update(id int, assessmentRequest models.AssessmentRequests, profile models.Profile) (body models.AssessmentResponse, err error) {
	assessmentData := assessmentRequest.Assessment_data
	switch assessmentRequest.Assessment_type {
	case "project":
		jsonData, _ := json.Marshal(assessmentData)
		var project models.AssessmentProject
		r := u.db.Table("assessment_project").Where("profile_id = ?", id).Preload("profile").First(&project)
		if r.RowsAffected == 0 {
			return body, fmt.Errorf("no data for assessment project")
		}
		json.Unmarshal(jsonData, &project)
		r = u.db.Session(&gorm.Session{FullSaveAssociations: true}).Table("assessment_project").Where("id = ?", project.Id).Select("*").Updates(&project)
		if err := r.Error; err != nil {
			return body, err
		}
		assessmentData = project
	case "progress":
		jsonData, _ := json.Marshal(assessmentData)
		var progress models.AssessmentProgress
		progress = models.AssessmentProgress{
			Progress_status: true,
		}
		r := u.db.Table("assessment_progress").Where("profile_id = ?", id).Preload("profile").First(&progress)
		if r.RowsAffected == 0 {
			return body, fmt.Errorf("no data for assessment progress")
		}
		json.Unmarshal(jsonData, &progress)
		r = u.db.Session(&gorm.Session{FullSaveAssociations: true}).Table("assessment_progress").Where("id = ?", progress.Id).Select("*").Updates(&progress)
		if err := r.Error; err != nil {
			return body, err
		}
		assessmentData = progress
	case "report":
		jsonData, _ := json.Marshal(assessmentData)
		var report models.AssessmentReport
		report = models.AssessmentReport{
			Report_status: true,
		}
		r := u.db.Table("assessment_project").Where("profile_id = ?", id).Preload("profile").First(&report)
		if r.RowsAffected == 0 {
			return body, fmt.Errorf("no data for assessment report")
		}
		json.Unmarshal(jsonData, &report)
		r = u.db.Session(&gorm.Session{FullSaveAssociations: true}).Table("assessment_report").Where("id = ?", report.Id).Select("*").Updates(&report)
		if err := r.Error; err != nil {
			return body, err
		}
		assessmentData = report
	case "article":
		jsonData, _ := json.Marshal(assessmentData)
		var article models.AssessmentArticle
		article = models.AssessmentArticle{
			Article_status: true,
		}
		r := u.db.Table("assessment_article").Where("profile_id = ?", id).Preload("profile").First(&article)
		if r.RowsAffected == 0 {
			return body, fmt.Errorf("no data for assessment article")
		}
		json.Unmarshal(jsonData, &article)
		r = u.db.Session(&gorm.Session{FullSaveAssociations: true}).Table("assessment_article").Where("id = ?", article.Id).Select("*").Updates(&article)
		if err := r.Error; err != nil {
			return body, err
		}
		assessmentData = article
	default:
		return body, fmt.Errorf("err")
	}

	body = models.AssessmentResponse{
		ProfileID:       assessmentRequest.ProfileID,
		Assessment_type: assessmentRequest.Assessment_type,
		Assessment_data: assessmentData,
	}
	return body, err

}
func (u *AssessmentHandler) create(assessmentRequest models.AssessmentRequests) (body models.AssessmentResponse, err error) {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	object := assessmentRequest.Assessment_data
	switch assessmentRequest.Assessment_type {
	case "project":
		jsonData, _ := json.Marshal(object)
		var project models.AssessmentProject
		project = models.AssessmentProject{
			Profile_id:     assessmentRequest.ProfileID,
			Project_status: true, // Set project_status to true
		}
		json.Unmarshal(jsonData, &project)
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&project).Error; err != nil {
			tx.Rollback()
			return body, err
		}
		object = &project
	case "progress":
		jsonData, _ := json.Marshal(object)
		var progress models.AssessmentProgress
		progress = models.AssessmentProgress{
			Profile_id:      assessmentRequest.ProfileID,
			Progress_status: true, // Set progress_status to true
		}
		json.Unmarshal(jsonData, &progress)
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&progress).Error; err != nil {
			tx.Rollback()
			return body, err
		}
		object = &progress
	case "report":
		jsonData, _ := json.Marshal(object)
		var report models.AssessmentReport
		report = models.AssessmentReport{
			Profile_id:    assessmentRequest.ProfileID,
			Report_status: true, // Set report_status to true
		}
		json.Unmarshal(jsonData, &report)
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&report).Error; err != nil {
			tx.Rollback()
			return body, err
		}
		object = &report
	case "article":
		jsonData, _ := json.Marshal(object)
		var article models.AssessmentArticle
		article = models.AssessmentArticle{
			Profile_id:     assessmentRequest.ProfileID,
			Article_status: true, // Set article_status to true
		}
		json.Unmarshal(jsonData, &article)
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&article).Error; err != nil {
			tx.Rollback()
			return body, err
		}
		object = &article
	default:
		return body, fmt.Errorf("")
	}
	tx.Commit()
	body = models.AssessmentResponse{
		ProfileID:       assessmentRequest.ProfileID,
		Assessment_type: assessmentRequest.Assessment_type,
		Assessment_data: object,
	}

	return body, err
}
