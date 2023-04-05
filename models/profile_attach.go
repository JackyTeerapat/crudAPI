package models

import "time"

type Profile_attach struct {
	Id           int       `json:"id"`
	File_name    string    `json:"file_name"`
	File_action  string    `json:"file_action"`
	File_storage string    `json:"file_storage"`
	Profile_id   int       `json:"profile_id"`
	Created_by   string    `json:"created_by"`
	Activated    bool      `json:"activated" gorm:"default:true"`
	Updated_by   string    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Profile_attach) TableName() string {
	return "profile_attach"
}
