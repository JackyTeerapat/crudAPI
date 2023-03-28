package models

import "time"

type Position struct {
	ID            int       `json:"id"`
	Position_name string    `json:"position"`
	Created_by    string    `json:"created_by"`
	Updated_by    string    `json:"updated_by"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Position) TableName() string {
	return "position"
}
