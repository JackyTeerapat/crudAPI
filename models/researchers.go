package models

import "time"

type Researcher struct {
	ProfileID   int       `json:"profile_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	University  string    `json:"university"`
	AddressHome string    `json:"address_home"`
	AddressWork string    `json:"address_work"`
	Email       string    `json:"email"`
	Degrees     []Degree  `json:"degrees"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
type Researchers []Researcher
