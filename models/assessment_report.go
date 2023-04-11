package models

import "time"

type AssessmentReport struct {
	Id               int       `json:"report_id"`
	Report_year      string    `json:"report_year"`
	Report_title     string    `json:"report_title"`
	Report_estimate  bool      `json:"report_estimate"`
	Estimate_remark  string    `json:"-"`
	Report_recommend bool      `json:"report_recommend"`
	Recommend_remark string    `json:"-"`
	File_name        string    `json:"file_name"`
	File_Id          int       `json:"file_id"`
	File_storage     string    `json:"-"`
	Period           bool      `json:"period"`
	Created_by       string    `json:"-"`
	Updated_by       string    `json:"-"`
	CreatedAt        time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *AssessmentReport) TableName() string {
	return "assessment_report"
}
