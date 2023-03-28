package models

type assessment_report struct {
	Id               int    `json:"id"`
	Report_year      string `json:"report_year"`
	Report_title     string `json:"report_title"`
	Report_estimate  int    `json:"report_estimate"`
	Estimate_remark  string `json:"estimate_remark"`
	Report_recommend int    `json:"report_recommend"`
	Recommend_remark string `json:"recommend_remark"`
	File_name        string `json:"file_name"`
	File_storage     string `json:"file_storage"`
	Period           int    `json:"period"`
	Created_by       string `json:"created_by"`
	Updated_by       string `json:"updated_by"`
	Created_at       string `json:"created_at"`
	Updated_at       string `json:"updated_at"`
}
