package minioclient

import (
	"CRUD-API/models"
	"time"

	"gorm.io/gorm"
)

func UpSertProfileAttach(db *gorm.DB, filename, data_type string, profile_id int) error {

	var profile_attach models.Profile_attach

	// check exist
	r := db.Table("profile_attach").Where("profile_id = ?", profile_id).Where("File_action = ?", data_type).First(&profile_attach)
	if r.RowsAffected == 0 {
		profile_attach.File_name = filename
		profile_attach.File_storage = data_type + "/" + filename
		profile_attach.File_action = data_type
		profile_attach.Activated = true
		profile_attach.Profile_id = profile_id
		r = db.Create(&profile_attach)
		if err := r.Error; err != nil {
			return err
		}
	} else {
		//update
		profile_attach.File_name = filename
		profile_attach.File_storage = data_type + "/" + filename
		profile_attach.Activated = true
		profile_attach.UpdatedAt = time.Now()
		r = db.Table("profile_attach").Updates(&profile_attach)
		if err := r.Error; err != nil {
			return err
		}
	}
	return nil
}

func UpdateAssessment(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var assessment models.Assessment
	res_flie_name := ""
	// check exist
	r := db.Table("assessment").Where("id = ?", profile_id).First(&assessment)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = assessment.Assessment_file_name
			assessment.Assessment_file_name = ""
			assessment.Assessment_file_action = ""
		} else {
			assessment.Assessment_file_name = filename
			assessment.Assessment_file_action = data_type
		}
		assessment.UpdatedAt = time.Now()
		r = db.Table("assessment").Select("assessment_file_name", "assessment_file_action", "updated_at").Updates(&assessment)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func UpdateAssessmentProject(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var project models.AssessmentProject
	res_flie_name := ""
	// check exist
	r := db.Table("assessment_project").Where("id = ?", profile_id).First(&project)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = project.File_name
			project.File_name = ""
			project.File_action = ""
		} else {
			project.File_name = filename
			project.File_action = data_type
		}
		project.UpdatedAt = time.Now()
		r = db.Table("assessment_project").Select("file_name", "file_action", "updated_at").Updates(&project)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func UpdateAssessmentProgress(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var progress models.AssessmentProgress
	res_flie_name := ""

	// check exist
	r := db.Table("assessment_progress").Where("id = ?", profile_id).First(&progress)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = progress.File_name
			progress.File_name = ""
			progress.File_action = ""
		} else {
			progress.File_name = filename
			progress.File_action = data_type
		}
		progress.UpdatedAt = time.Now()
		r = db.Table("assessment_progress").Select("file_name", "file_action", "updated_at").Updates(&progress)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func UpdateAssessmentReport(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var repport models.AssessmentReport
	res_flie_name := ""

	// check exist
	r := db.Table("assessment_report").Where("id = ?", profile_id).First(&repport)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = repport.File_name
			repport.File_name = ""
			repport.File_action = ""
		} else {
			repport.File_name = filename
			repport.File_action = data_type
		}
		repport.UpdatedAt = time.Now()
		r = db.Table("assessment_report").Select("file_name", "file_action", "updated_at").Updates(&repport)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func UpdateAssessmentArticle(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var article models.AssessmentArticle
	res_flie_name := ""

	// check exist
	r := db.Table("assessment_article").Where("id = ?", profile_id).First(&article)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = article.File_name
			article.File_name = ""
			article.File_action = ""
		} else {
			article.File_name = filename
			article.File_action = data_type
		}
		article.UpdatedAt = time.Now()
		r = db.Table("assessment_article").Select("file_name", "file_action", "updated_at").Updates(&article)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func DeleteProfileAttach(db *gorm.DB, data_type string, profile_id int) (string, error) {

	var profile_attach models.Profile_attach
	res_flie_name := ""
	// check exist
	r := db.Table("profile_attach").Where("profile_id = ?", profile_id).Where("File_action = ?", data_type).First(&profile_attach)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		res_flie_name = profile_attach.File_name
		profile_attach.File_name = ""
		profile_attach.File_storage = ""
		profile_attach.Activated = false
		profile_attach.UpdatedAt = time.Now()
		r = db.Table("profile_attach").Select("file_name", "file_storage", "activated", "updated_at").Updates(&profile_attach)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func GetAssessment(db *gorm.DB, data_type string, profile_id int) (string, error) {

	var assessment models.Assessment
	res_flie_name := ""
	// check exist
	r := db.Table("assessment").Where("id = ?", profile_id).First(&assessment)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	}
	res_flie_name = assessment.Assessment_file_name
	return res_flie_name, nil
}

func GetAssessmentProject(db *gorm.DB, data_type string, profile_id int) (string, error) {

	var project models.AssessmentProject
	res_flie_name := ""
	// check exist
	r := db.Table("assessment_project").Where("id = ?", profile_id).First(&project)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	}
	res_flie_name = project.File_name
	return res_flie_name, nil
}

func GetAssessmentProgress(db *gorm.DB, data_type string, profile_id int) (string, error) {

	var progress models.AssessmentProgress
	res_flie_name := ""

	// check exist
	r := db.Table("assessment_progress").Where("id = ?", profile_id).First(&progress)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	}
	res_flie_name = progress.File_name
	return res_flie_name, nil
}

func GetAssessmentReport(db *gorm.DB, data_type string, profile_id int) (string, error) {

	var repport models.AssessmentReport
	res_flie_name := ""

	// check exist
	r := db.Table("assessment_report").Where("id = ?", profile_id).First(&repport)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	}
	res_flie_name = repport.File_name
	return res_flie_name, nil
}

func GetAssessmentArticle(db *gorm.DB, data_type string, profile_id int) (string, error) {

	var article models.AssessmentArticle
	res_flie_name := ""

	// check exist
	r := db.Table("assessment_article").Where("id = ?", profile_id).First(&article)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	}
	res_flie_name = article.File_name
	return res_flie_name, nil
}

func GetProfileAttach(db *gorm.DB, data_type string, profile_id int) (string, error) {

	var profile_attach models.Profile_attach
	res_flie_name := ""
	// check exist
	r := db.Table("profile_attach").Where("profile_id = ?", profile_id).Where("File_action = ?", data_type).First(&profile_attach)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	}
	res_flie_name = profile_attach.File_name
	return res_flie_name, nil
}

func DeleteProfile(db *gorm.DB, profile_id int) error {
	r := db.Delete(&models.Profile{}, profile_id)
	if err := r.Error; err != nil {
		return err
	}
	return nil
}

func DeleteAssessment(db *gorm.DB, assessment_id int) error {
	r := db.Delete(&models.Assessment{}, assessment_id)
	if err := r.Error; err != nil {
		return err
	}
	return nil
}
