package models

import "time"

type Exploration struct {
	Id             int       `json:"id"`
	Explore_name   string    `json:"explore_name"`
	Explore_year   string    `json:"explore_year"`
	Explore_detail string    `json:"explore_detail"`
	Profile_id     int       `json:"profile_id"`
	Created_by     string    `json:"created_by"`
	Updated_by     string    `json:"updated_by"`
	Activated      bool      `json:"activated" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Exploration) TableName() string {
	return "exploration"
}
