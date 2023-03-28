package models

type Profile struct {
	Id           int    `json:"id"`
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	Position_id  int    `json:"position_id"`
	University   string `json:"university_id"`
	Address_home string `json:"address_home"`
	Address_work string `json:"address_work"`
	Email        string `json:"email"`
	Phone_number string `json:"phone_number"`
	Created_by   string `json:"created_by"`
	Updated_by   string `json:"updated_by"`
	Created_at   string `json:"created_at"`
	Updated_at   string `json:"updated_at"`
}
