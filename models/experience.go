package models

import "time"

type Experience struct {
	Id                    int       `json:"id"`
	Experience_type       string    `json:"experience_type"`
	Experience_start      string    `json:"experience_start"`
	Experience_end        string    `json:"experience_end"`
	Experience_university string    `json:"experience_university"`
	Experience_remark     string    `json:"experience_remark"`
	Profile_id            int       `json:"profile_id"`
	Created_by            string    `json:"created_by"`
	Updated_by            string    `json:"updated_by"`
	CreatedAt             time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Experience) TableName() string {
	return "experience"
}
