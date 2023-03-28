package models

type Position struct {
	ID            int    `json:"id"`
	Position_name string `json:"position"`
	Created_by    string `json:"created_by"`
	Updated_by    string `json:"updated_by"`
	Created_at    string `json:"created_at"`
	Updated_at    string `json:"updated_at"`
}
