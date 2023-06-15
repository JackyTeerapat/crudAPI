package models

import "time"

type AssessmentArticle struct {
	Id                int       `json:"article_id"`
	Profile_id        int       `json:"profile_id"`
	Article_year      string    `json:"article_year"`
	Article_title     string    `json:"article_title"`
	Article_estimate  bool      `json:"article_estimate"`
	Article_recommend bool      `json:"article_recommend"`
	Article_type      string    `json:"article_type"`
	Article_status    string    `json:"article_status"`
	File_name         string    `json:"file_name"`
	File_action       string    `json:"file_action"`
	File_storage      string    `json:"-"`
	Period            bool      `json:"period"`
	Created_by        string    `json:"-"`
	Updated_by        string    `json:"-"`
	CreatedAt         time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"-" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
type AssessmentArticleGet struct {
	Id                int       `json:"article_id"`
	Profile_id        int       `json:"-"`
	Article_year      string    `json:"article_year"`
	Article_title     string    `json:"article_title"`
	Article_estimate  bool      `json:"article_estimate"`
	Article_recommend bool      `json:"article_recommend"`
	Article_type      string    `json:"article_type"`
	Article_status    string    `json:"article_status"`
	File_name         string    `json:"file_name"`
	File_action       string    `json:"file_action"`
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
