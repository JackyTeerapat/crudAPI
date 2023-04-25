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
		Where("profile_id = ?", id).
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

func (u *AssessmentHandler) CreateAssessmentHandler(c *gin.Context) {
	var assessment models.AssessmentRequests
	if err := c.ShouldBindJSON(&assessment); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	var checkAs models.Assessment
	ckeck := u.db.Table("assessment").Where("profile_id = ?", assessment.ProfileID).First(&checkAs)
	if ckeck.RowsAffected != 0 {
		e := u.RollbackDeleteProFile(assessment.ProfileID)
		if e != nil {
			res := api.ResponseApi(http.StatusBadRequest, nil, e)
			c.JSON(http.StatusBadRequest, res)
			return
		}
		res := api.ResponseApi(http.StatusBadRequest, nil, fmt.Errorf("This profile has already create assessment"))
		c.JSON(http.StatusBadRequest, res)
		return
	}
	body, err := u.create(assessment)
	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	var resAssessment models.Assessment
	r := u.db.Table("assessment").
		Where("id = ?", body.Id).
		Preload("Project").
		Preload("Progress").
		Preload("Report").
		Preload("Article").
		First(&resAssessment)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}
	if err := r.Error; err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res := api.ResponseApi(http.StatusCreated, resAssessment, nil)
	c.JSON(http.StatusCreated, res)
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
		Id:               assessment.Id,
		Assessment_start: assessmentRequest.AssessmentStart,
		Assessment_end:   assessmentRequest.AssessmentEnd,
		Project: models.AssessmentProject{
			Id:                assessment.ProjectID,
			Project_year:      assessmentRequest.Project_year,
			Project_title:     assessmentRequest.Project_title,
			Project_estimate:  assessmentRequest.Project_estimate,
			Project_recommend: assessmentRequest.Project_recommend,
			Period:            assessmentRequest.Project_period,
			Updated_by:        "admin",
		},
		Progress: models.AssessmentProgress{
			Id:                 assessment.ProgressID,
			Progress_year:      assessmentRequest.Progress_year,
			Progress_title:     assessmentRequest.Progress_title,
			Progress_estimate:  assessmentRequest.Progress_estimate,
			Progress_recommend: assessmentRequest.Progress_recommend,
			Period:             assessmentRequest.Progress_period,
			Updated_by:         "admin",
		},
		Report: models.AssessmentReport{
			Id:               assessment.ReportID,
			Report_year:      assessmentRequest.Report_year,
			Report_title:     assessmentRequest.Report_title,
			Report_estimate:  assessmentRequest.Report_estimate,
			Report_recommend: assessmentRequest.Report_recommend,
			Period:           assessmentRequest.Report_period,
			Updated_by:       "admin",
		},
		Article: models.AssessmentArticle{
			Id:                assessment.ArticleID,
			Article_year:      assessmentRequest.Article_year,
			Article_title:     assessmentRequest.Article_title,
			Article_estimate:  assessmentRequest.Article_estimate,
			Article_recommend: assessmentRequest.Article_recommend,
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
	project := models.AssessmentProject{
		Project_year:      assessmentRequest.Project_year,
		Project_title:     assessmentRequest.Project_title,
		Project_point:     assessmentRequest.Project_point,
		Project_estimate:  assessmentRequest.Project_estimate,
		Project_recommend: assessmentRequest.Project_recommend,
		Period:            assessmentRequest.Project_period,
	}
	p := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&project)
	if err := p.Error; err != nil {
		// e := u.RollbackDeleteProFile(assessmentRequest.ProfileID)
		// if e != nil {
		// 	return body, e
		// }
		tx.Rollback()
		return body, err
	}

	progress := models.AssessmentProgress{
		Progress_year:      assessmentRequest.Progress_year,
		Progress_title:     assessmentRequest.Progress_title,
		Progress_estimate:  assessmentRequest.Progress_estimate,
		Progress_recommend: assessmentRequest.Project_recommend,
		Period:             assessmentRequest.Progress_period,
	}

	if err := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&progress).Error; err != nil {
		// e := u.RollbackDeleteProFile(assessmentRequest.ProfileID)
		// if e != nil {
		// 	return body, e
		// }
		tx.Rollback()
		return body, err
	}
	report := models.AssessmentReport{
		Report_year:      assessmentRequest.Report_year,
		Report_title:     assessmentRequest.Report_title,
		Report_estimate:  assessmentRequest.Report_estimate,
		Report_recommend: assessmentRequest.Report_recommend,
		Period:           assessmentRequest.Report_period,
	}

	if err := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&report).Error; err != nil {
		// e := u.RollbackDeleteProFile(assessmentRequest.ProfileID)
		// if e != nil {
		// 	return body, e
		// }
		tx.Rollback()
		return body, err
	}
	article := models.AssessmentArticle{
		Article_year:      assessmentRequest.Article_year,
		Article_title:     assessmentRequest.Article_title,
		Article_estimate:  assessmentRequest.Article_estimate,
		Article_recommend: assessmentRequest.Article_recommend,
		Period:            assessmentRequest.Article_period,
	}

	if err := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&article).Error; err != nil {
		// e := u.RollbackDeleteProFile(assessmentRequest.ProfileID)
		// if e != nil {
		// 	return body, e
		// }
		tx.Rollback()
		return body, err
	}
	body = models.Assessment{
		Assessment_start:       assessmentRequest.AssessmentStart,
		Assessment_end:         assessmentRequest.AssessmentEnd,
		Assessment_file_action: "assessment",
		Assessment_file_name:   "test",
		ProjectID:              project.Id,
		ProgressID:             progress.Id,
		ProfileID:              assessmentRequest.ProfileID,
		ReportID:               report.Id,
		ArticleID:              article.Id,
		Created_by:             "admin",
		Updated_by:             "admin",
	}

	r := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&body)
	if err := r.Error; err != nil {
		tx.Rollback()
		return body, err
	}
	return body, err
}

func (u *AssessmentHandler) RollbackDeleteProFile(profileId int) (err error) {
	dassessment := u.db.Where("profile_id = ?", profileId).Delete(&models.Assessment{})
	if err := dassessment.Error; err != nil {
		return err
	}
	dprogram := u.db.Where("profile_id = ?", profileId).Delete(&models.Program{})
	if err = dprogram.Error; err != nil {
		return err
	}
	dProfileAttach := u.db.Where("profile_id = ?", profileId).Delete(&models.Profile_attach{})
	if err = dProfileAttach.Error; err != nil {
		return err
	}
	dexploration := u.db.Where("profile_id = ?", profileId).Delete(&models.Exploration{})
	if err = dexploration.Error; err != nil {
		return err
	}
	dexperience := u.db.Where("profile_id = ?", profileId).Delete(&models.Experience{})
	if err = dexperience.Error; err != nil {
		return err
	}
	ddegree := u.db.Where("profile_id = ?", profileId).Delete(&models.Degree{})
	if err = ddegree.Error; err != nil {
		return err
	}
	dprofile := u.db.Delete(&models.Profile{}, profileId)
	if err = dprofile.Error; err != nil {
		return err
	}

	return nil
}
