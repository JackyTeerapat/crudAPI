package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (u *User) TableName() string {
	return "users"
}
