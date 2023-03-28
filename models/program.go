package models

import "time"

type Program struct {
	Id           int       `json:"id"`
	Program_name string    `json:"program_name"`
	Profile_id   int       `json:"profile_id"`
	Created_by   string    `json:"created_by"`
	Updated_by   string    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Program) TableName() string {
	return "program"
}