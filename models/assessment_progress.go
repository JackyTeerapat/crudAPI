package models

type assessment_progress struct {
	Id                 int    `json:"id"`
	Progress_year      string `json:"progress_year"`
	Progress_title     string `json:"progress_title"`
	Progress_estimate  int    `json:"progress_estimate"`
	Estimate_remark    string `json:"estimate_remark"`
	Progress_recommend int    `json:"progress_recommend"`
	Recommend_remark   string `json:"recommend_remark"`
	File_name          string `json:"file_name"`
	File_storage       string `json:"file_storage"`
	Period             int    `json:"period"`
	Created_by         string `json:"created_by"`
	Updated_by         string `json:"updated_by"`
	Created_at         string `json:"created_at"`
	Updated_at         string `json:"updated_at"`
}
