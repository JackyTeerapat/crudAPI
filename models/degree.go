package models

type degrees struct {
	Id                int    `json:"id"`
	Degree_type       string `json:"degree_type"`
	Degree_program    string `json:"degree_program"`
	Degree_university string `json:"degree_university"`
	Profile_id        int    `json:"profile_id"`
	Created_by        string `json:"created_by"`
	Updated_by        string `json:"updated_by"`
	Created_at        string `json:"created_at"`
	Updated_at        string `json:"updated_at"`
}
