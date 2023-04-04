package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `json:"id"`
	Username string `json:"username" binding:"required" gorm:"unique"` // ชื่อผู้ใช้งาน
	Role     string `json:"role" binding:"required"`                   // สิทธิ์ผู้ใช้งาน
}

type Register struct {
	ID       uint   `json:"id"`
	Username string `json:"username" binding:"required" gorm:"unique"`
	Role     string `json:"role" binding:"required"`
}

type LoginRespones struct {
	ID       uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

func (u *User) TableName() string {
	return "users"
}
