package models

type Program struct {
	Id           int    `json:"id"`
	Program_name string `json:"program_name"`
	Profile_id   int    `json:"profile_id"`
	Created_by   string `json:"created_by"`
	Updated_by   string `json:"updated_by"`
	Created_at   string `json:"created_at"`
	Updated_at   string `json:"updated_at"`
}
