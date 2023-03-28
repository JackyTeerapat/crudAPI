package models

type assessment_project struct {
	Id                int    `json:"id"`
	Project_year      string `json:"project_year"`
	Project_title     string `json:"project_title"`
	Project_point     int    `json:"project_point"`
	Project_estimate  int    `json:"project_estimate"`
	Estimate_remark   string `json:"estimate_remark"`
	Project_recommend int    `json:"project_recommend"`
	Recommend_remark  string `json:"recommend_remark"`
	File_name         string `json:"file_name"`
	File_storage      string `json:"file_storage"`
	Period            int    `json:"period"`
	Created_by        string `json:"created_by"`
	Updated_by        string `json:"updated_by"`
	Created_at        string `json:"created_at"`
	Updated_at        string `json:"updated_at"`
}
