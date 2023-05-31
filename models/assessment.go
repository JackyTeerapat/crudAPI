package models

import (
	"time"
)

type Assessment struct {
	ProfileID               int                  `json:"profile_id"`
	Id                      int                  `json:"assessment_id"`
	Assessment_start        string               `json:"assessment_start"`
	Assessment_end          string               `json:"assessment_end"`
	Assessment_file_name    string               `json:"assessment_file_name"`
	Assessment_file_action  string               `json:"assessment_file_action"`
	Assessment_file_storage string               `json:"-"`
	ProjectID               int                  `json:"-"`
	Project                 []AssessmentProject  `gorm:"foreignkey:ProjectID"`
	ProgressID              int                  `json:"-"`
	Progress                []AssessmentProgress `gorm:"foreignkey:ProgressID"`
	ReportID                int                  `json:"-"`
	Report                  []AssessmentReport   `gorm:"foreignkey:ReportID"`
	ArticleID               int                  `json:"-"`
	Article                 []AssessmentArticle  `gorm:"foreignkey:ArticleID"`
	Created_by              string               `json:"-"`
	Updated_by              string               `json:"-"`
	CreatedAt               time.Time            `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt               time.Time            `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type AssessmentRequests struct {
	ProfileID       int         `json:"profile_id"`
	Assessment_type string      `json:"assessment_type"`
	Assessment_data interface{} `json:"assessment_data"`
}
type AssessmentResponse struct {
	ProfileID       int         `json:"profile_id"`
	Assessment_type string      `json:"assessment_type"`
	Assessment_data interface{} `json:"assessment_data"`
}

func (p *Assessment) TableName() string {
	return "assessment"
}
