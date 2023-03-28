package models

type assessment struct {
	Id                      int    `json:"id"`
	Assessment_start        string `json:"assessment_start"`
	Assessment_end          string `json:"assessment_end"`
	Assessment_file_name    string `json:"assessment_file_name"`
	Assessment_file_storage string `json:"assessment_file_storage"`
	Project_id              int    `json:"project_id"`
	Progress_id             int    `json:"progress_id"`
	Report_id               int    `json:"report_id"`
	Article_id              int    `json:"article_id"`
	Profile_id              int    `json:"profile_id"`
	Created_by              string `json:"created_by"`
	Updated_by              string `json:"updated_by"`
	Created_at              string `json:"created_at"`
	Updated_at              string `json:"updated_at"`
}
