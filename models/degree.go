package models

import "time"

type Degree struct {
	ID                int       `json:"id"`
	Degree_type       string    `json:"degree_type"`
	Degree_program    string    `json:"degree_program"`
	Degree_university string    `json:"degree_university"`
	Profile_id        int       `json:"profile_id"`
	Created_by        string    `json:"created_by"`
	Updated_by        string    `json:"updated_by"`
	Activated         bool      `json:"activated" gorm:"default:true"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Degree) TableName() string {
	return "degree"
}
