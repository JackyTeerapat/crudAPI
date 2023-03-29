package models

import "time"

type Article struct {
	Id                int       `json:"id"`
	Article_year      string    `json:"article_year"`
	Article_title     string    `json:"article_title"`
	Article_estimate  bool      `json:"article_estimate"`
	Estimate_remark   string    `json:"estimate_remark"`
	Article_recommend bool      `json:"article_recommend"`
	Recommend_remark  string    `json:"recommend_remark"`
	File_name         string    `json:"file_name"`
	File_storage      string    `json:"file_storage"`
	Period            bool      `json:"period"`
	Created_by        string    `json:"created_by"`
	Updated_by        string    `json:"updated_by"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *Article) TableName() string {
	return "assessment_article"
}
