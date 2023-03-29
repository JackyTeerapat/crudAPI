package models

import "time"

type Assessment_project struct {
	Id                int       `json:"id"`
	Project_year      string    `json:"project_year"`
	Project_title     string    `json:"project_title"`
	Project_point     int       `json:"project_point"`
	Project_estimate  bool      `json:"project_estimate"`
	Estimate_remark   string    `json:"estimate_remark"`
	Project_recommend bool      `json:"project_recommend"`
	Recommend_remark  string    `json:"recommend_remark"`
	File_name         string    `json:"file_name"`
	File_storage      string    `json:"file_storage"`
	Period            bool      `json:"period"`
	Created_by        string    `json:"created_by"`
	Updated_by        string    `json:"updated_by"`
	Created_at        time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Updated_at        time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Assessment_project) TableName() string {
	return "assessment_project"
}
