package models

import "time"

type AssessmentArticle struct {
	Id                int       `json:"article_id"`
	Article_year      string    `json:"article_year"`
	Article_title     string    `json:"article_title"`
	Article_estimate  bool      `json:"article_estimate"`
	Estimate_remark   string    `json:"-"`
	Article_recommend bool      `json:"article_recommend"`
	Recommend_remark  string    `json:"-"`
	File_name         string    `json:"file_name"`
	File_Id           int       `json:"file_id"`
	File_storage      string    `json:"-"`
	Period            bool      `json:"period"`
	Created_by        string    `json:"-"`
	Updated_by        string    `json:"-"`
	CreatedAt         time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (p *AssessmentArticle) TableName() string {
	return "assessment_article"
}
