package models

import "time"

type Progress struct {
	Id                 int       `json:"id"`
	Progress_year      string    `json:"progress_year"`
	Progress_title     string    `json:"progress_title"`
	Progress_estimate  bool      `json:"progress_estimate"`
	Estimate_remark    string    `json:"estimate_remark"`
	Progress_recommend bool      `json:"progress_recommend"`
	Recommend_remark   string    `json:"recommend_remark"`
	File_name          string    `json:"file_name"`
	File_storage       string    `json:"file_storage"`
	Period             bool      `json:"period"`
	Created_by         string    `json:"created_by"`
	Updated_by         string    `json:"updated_by"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Progress) TableName() string {
	return "assessment_progress"
}
