package models

import "time"

type Profile struct {
	ID           int       `json:"id"`
	First_name   string    `json:"first_name"`
	Last_name    string    `json:"last_name"`
	PositionID   int       `json:"position_id"`
	Position     Position  `gorm:"foreignkey:PositionID"`
	University   string    `json:"university"`
	Address_home string    `json:"address_home"`
	Address_work string    `json:"address_work"`
	Email        string    `json:"email"`
	Phone_number string    `json:"phone_number"`
	Created_by   string    `json:"created_by"`
	Updated_by   string    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Profile) TableName() string {
	return "profile"
}
