package models

import "time"

type Report struct {
	Id               int       `json:"id"`
	Report_year      string    `json:"report_year"`
	Report_title     string    `json:"report_title"`
	Report_estimate  bool      `json:"report_estimate"`
	Estimate_remark  string    `json:"estimate_remark"`
	Report_recommend bool      `json:"report_recommend"`
	Recommend_remark string    `json:"recommend_remark"`
	File_name        string    `json:"file_name"`
	File_storage     string    `json:"file_storage"`
	Period           bool      `json:"period"`
	Created_by       string    `json:"created_by"`
	Updated_by       string    `json:"updated_by"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Report) TableName() string {
	return "assessment_report"
}
