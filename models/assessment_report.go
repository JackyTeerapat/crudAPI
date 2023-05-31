package models

import "time"

type AssessmentReport struct {
	Id               int       `json:"report_id"`
	Profile_id       int       `json:"profile_id"`
	Report_year      string    `json:"report_year"`
	Report_funding   string    `json:"report_funding"`
	Report_source    string    `json:"report_source"`
	Report_title     string    `json:"report_title"`
	Report_estimate  bool      `json:"report_estimate"`
	Report_recommend bool      `json:"report_recommend"`
	File_name        string    `json:"file_name"`
	File_action      string    `json:"file_action"`
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
