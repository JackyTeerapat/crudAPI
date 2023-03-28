package models

type exploration struct {
	Id             int    `json:"id"`
	Explore_name   string `json:"explore_name"`
	Explore_year   string `json:"explore_year"`
	Explore_detail string `json:"explore_detail"`
	Profile_id     int    `json:"profile_id"`
	Created_by     string `json:"created_by"`
	Updated_by     string `json:"updated_by"`
	Created_at     string `json:"created_at"`
	Updated_at     string `json:"updated_at"`
}
