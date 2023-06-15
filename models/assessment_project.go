package models

import "time"

type AssessmentProject struct {
	Id                int       `json:"project_id"`
	Profile_id        int       `json:"profile_id"`
	Project_year      string    `json:"project_year"`
	Project_funding   string    `json:"project_funding"`
	Project_source    string    `json:"project_source"`
	Project_title     string    `json:"project_title"`
	Project_point     int       `json:"project_point"`
	Project_estimate  bool      `json:"project_estimate"`
	Project_recommend bool      `json:"project_recommend"`
	Project_creator   string    `json:"project_creator"`
	Project_status    bool      `json:"project_status"`
	File_name         string    `json:"file_name"`
	File_action       string    `json:"file_action"`
	File_storage      string    `json:"-"`
	Period            bool      `json:"period"`
	Created_by        string    `json:"-"`
	Updated_by        string    `json:"-"`
	CreatedAt         time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
type AssessmentProjectGet struct {
	Id                int       `json:"project_id"`
	Profile_id        int       `json:"-"`
	Project_year      string    `json:"project_year"`
	Project_funding   string    `json:"project_funding"`
	Project_source    string    `json:"project_source"`
	Project_title     string    `json:"project_title"`
	Project_point     int       `json:"project_point"`
	Project_estimate  bool      `json:"project_estimate"`
	Project_recommend bool      `json:"project_recommend"`
	Project_creator   string    `json:"project_creator"`
	Project_status    bool      `json:"project_status"`
	File_name         string    `json:"file_name"`
	File_action       string    `json:"file_action"`
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
