package models

import "gorm.io/gorm"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" binding:"required"` // ชื่อผู้ใช้งาน
	Role     string `json:"role" binding:"required"`     // สิทธิ์ผู้ใช้งาน
}

type Register struct {
	gorm.Model
	Username string `json:"username" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (u *User) TableName() string {
	return "users"
}
