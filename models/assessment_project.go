package models

import "time"

type AssessmentProject struct {
	Id                int       `json:"project_id"`
	Project_year      string    `json:"project_year"`
	Project_title     string    `json:"project_title"`
	Project_point     int       `json:"project_point"`
	Project_estimate  bool      `json:"project_estimate"`
	Project_recommend bool      `json:"project_recommend"`
	File_name         string    `json:"file_name"`
	File_Id           int       `json:"file_id"`
	File_storage      string    `json:"-"`
	Period            bool      `json:"period"`
	Created_by        string    `json:"-"`
	Updated_by        string    `json:"-"`
	CreatedAt         time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *AssessmentProject) TableName() string {
	return "assessment_project"
}
