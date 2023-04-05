package models

import "time"

type AssessmentProject struct {
	ID                int       `json:"id"`
	Project_year      string    `json:"project_year"`
	Project_title     string    `json:"project_title"`
	Project_point     int       `json:"project_point"`
	Project_estimate  bool      `json:"project_estimate"`
	Project_recommend bool      `json:"project_recommend"`
	File_name         string    `json:"file_name"`
	File_storage      string    `json:"file_storage"`
	Period            bool      `json:"period"`
	Created_by        string    `json:"created_by"`
	Updated_by        string    `json:"updated_by"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *AssessmentProject) TableName() string {
	return "assessment_project"
}
