package models

type assessment_article struct {
	Id                int    `json:"id"`
	Article           string `json:"article"`
	Article_title     string `json:"article_title"`
	Article_estimate  int    `json:"article_estimate"`
	Estimate_remark   string `json:"estimate_remark"`
	Article_recommend int    `json:"article_recommend"`
	Recommend_remark  string `json:"recommend_remark"`
	File_name         string `json:"file_name"`
	File_storage      string `json:"file_storage"`
	Period            int    `json:"period"`
	Created_by        string `json:"created_by"`
	Updated_by        string `json:"updated_by"`
	Created_at        string `json:"created_at"`
	Updated_at        string `json:"updated_at"`
}
