package models

import (
	"time"
)

type Assessment struct {
	ID                      uint `gorm:"primarykey"`
	Assessment_start        string
	Assessment_end          string
	Assessment_file_name    string
	Assessment_file_storage string
	ProjectID               int
	Project                 AssessmentProject `gorm:"foreignKey:ProjectID"`
	ProgressID              int
	Progress                AssessmentProgress `gorm:"foreignKey:ProgressID"`
	ReportID                int
	Report                  AssessmentReport `gorm:"foreignKey:ReportID"`
	ArticleID               int
	Article                 AssessmentArticle `gorm:"foreignKey:ArticleID"`
	ProfileID               int
	Profile                 Profile `gorm:"foreignKey:ProfileID"`
	Created_by              string
	Updated_by              string
	CreatedAt               time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt               time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type AssessmentRequests struct {
	Id                   int    `json:"id"`
	Assessment_start     string `json:"assessment_start"`
	Assessment_end       string `json:"assessment_end"`
	Assessment_file_name string `json:"assessment_file"`
	Project_year         string `json:"project_year"`
	Project_title        string `json:"project_title"`
	Project_point        int    `json:"project_point"`
	Project_estimate     bool   `json:"project_estimate"`
	Project_recommend    bool   `json:"project_recommend"`
	Project_file         string `json:"project_file"`
	Project_period       bool   `json:"project_period"`
	Progress_year        string `json:"progress_year"`
	Progress_title       string `json:"progress_title"`
	Progress_estimate    bool   `json:"progress_estimate"`
	Progress_recommend   bool   `json:"progress_recommend"`
	Progress_file        string `json:"progress_file"`
	Progress_period      bool   `json:"progress_period"`
	Report_year          string `json:"report_year"`
	Report_title         string `json:"report_title"`
	Report_estimate      bool   `json:"report_estimate"`
	Report_recommend     bool   `json:"report_recommend"`
	Report_file          string `json:"report_file"`
	Report_period        bool   `json:"report_period"`
	Article_year         string `json:"Article_year"`
	Article_title        string `json:"Article_title"`
	Article_estimate     bool   `json:"Article_estimate"`
	Article_recommend    bool   `json:"Article_recommend"`
	Article_file         string `json:"Article_file"`
	Article_period       bool   `json:"Article_period"`
}

func (p *Assessment) TableName() string {
	return "assessment"
}
