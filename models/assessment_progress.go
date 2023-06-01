package models

import "time"

type AssessmentProgress struct {
	Id                 int       `json:"progress_id"`
	Profile_id         int       `json:"profile_id"`
	Progress_year      string    `json:"progress_year"`
	Progress_funding   string    `json:"progress_funding"`
	Progress_source    string    `json:"progress_source"`
	Progress_title     string    `json:"progress_title"`
	Progress_estimate  bool      `json:"progress_estimate"`
	Progress_recommend bool      `json:"progress_recommend"`
	Progress_creator   string    `json:"progress_creator"`
	Progress_type      string    `json:"progress_type"`
	File_name          string    `json:"file_name"`
	File_action        string    `json:"file_action"`
	File_storage       string    `json:"-"`
	Period             bool      `json:"period"`
	Created_by         string    `json:"-"`
	Updated_by         string    `json:"-"`
	CreatedAt          time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type AssessmentProgressGet struct {
	Id                 int       `json:"progress_id"`
	Profile_id         int       `json:"-"`
	Progress_year      string    `json:"progress_year"`
	Progress_funding   string    `json:"progress_funding"`
	Progress_source    string    `json:"progress_source"`
	Progress_title     string    `json:"progress_title"`
	Progress_estimate  bool      `json:"progress_estimate"`
	Progress_recommend bool      `json:"progress_recommend"`
	Progress_creator   string    `json:"progress_creator"`
	Progress_type      string    `json:"progress_type"`
	File_name          string    `json:"file_name"`
	File_action        string    `json:"file_action"`
	File_storage       string    `json:"-"`
	Period             bool      `json:"period"`
	Created_by         string    `json:"-"`
	Updated_by         string    `json:"-"`
	CreatedAt          time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *AssessmentProgress) TableName() string {
	return "assessment_progress"
}
