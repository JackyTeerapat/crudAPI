package models

import "time"

type AssessmentProgress struct {
	Id                 int    `json:"progress_id"`
	Progress_year      string `json:"progress_year"`
	Progress_title     string `json:"progress_title"`
	Progress_estimate  bool   `json:"progress_estimate"`
	Progress_recommend bool   `json:"progress_recommend"`
	File_name          string `json:"file_name"`
	File_action        string `json:"file_action"`
	// File_Id            int       `json:"-"`
	File_storage string    `json:"-"`
	Period       bool      `json:"period"`
	Created_by   string    `json:"-"`
	Updated_by   string    `json:"-"`
	CreatedAt    time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *AssessmentProgress) TableName() string {
	return "assessment_progress"
}
