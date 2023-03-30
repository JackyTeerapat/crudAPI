package models

import "time"

type Assessment struct {
	Id                      int       `json:"id"`
	Assessment_start        string    `json:"assessment_start"`
	Assessment_end          string    `json:"assessment_end"`
	Assessment_file_name    string    `json:"assessment_file_name"`
	Assessment_file_storage string    `json:"assessment_file_storage"`
	ProjectID               int       `json:"project_id"`
	Project                 Project   `gorm:"foreignkey:ProjectID"`
	ProgressID              int       `json:"progress_id"`
	Progress                Progress  `gorm:"foreignkey:ProgressID"`
	ReportID                int       `json:"report_id"`
	Report                  Report    `gorm:"foreignkey:ReportID"`
	ArticleID               int       `json:"article_id"`
	Article                 Article   `gorm:"foreignkey:ArticleID"`
	ProfileID               int       `json:"profile_id"`
	Profile                 Profile   `gorm:"foreignkey:ProfileID"`
	Created_by              string    `json:"created_by"`
	Updated_by              string    `json:"updated_by"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt               time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Assessment) TableName() string {
	return "assessment"
}
