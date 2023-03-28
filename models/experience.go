package models

type experience struct {
	Id                    int    `json:"id"`
	Experience_type       string `json:"experience_type"`
	Experience_start      string `json:"experience_start"`
	Experience_end        string `json:"experience_end"`
	Experience_university string `json:"experience_university"`
	Experience_remark     string `json:"experience_remark"`
	Profile_id            int    `json:"profile_id"`
	Created_by            string `json:"created_by"`
	Updated_by            string `json:"updated_by"`
	Created_at            string `json:"created_at"`
	Updated_at            string `json:"updated_at"`
}
