package models

type profile_attach struct {
	Id           int    `json:"id"`
	File_name    string `json:"file_name"`
	File_action  string `json:"file_action"`
	File_storage string `json:"file_storage"`
	Profile_id   int    `json:"profile_id"`
	Created_by   string `json:"created_by"`
	Updated_by   string `json:"updated_by"`
	Created_at   string `json:"created_at"`
	Updated_at   string `json:"updated_at"`
}
