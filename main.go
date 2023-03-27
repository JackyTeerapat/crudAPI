package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

var users = []User{
	{ID: "1", Username: "champInwZa007", Role: "Admin"},
}

var dsn = "postgres://navjsbdt:CXbvdzgydzdeZKUi_WYzMxzxAjJqnYbF@satao.db.elephantsql.com/navjsbdt"
var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

func main() {
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&users)
	r := gin.New()

	r.GET("/user", listUser)
	r.POST("/user", createUserHandler)
	r.DELETE("/user/:id", deleteUserHandler)
	r.PUT("/user/:id", updateUserHandler)

	r.Run()
}

func listUser(c *gin.Context) {
	r := db.Find(&users)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func createUserHandler(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := db.Create(&user)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func deleteUserHandler(c *gin.Context) {
	id := c.Param("id")

	//รับมาแล้วสร้างเป็น ข้อมูล ลง Tables
	r := db.Delete(&User{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func updateUserHandler(c *gin.Context) {
	id := c.Param("id")
	var updatedUser User

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user User
	r := db.First(&user, "id = ?", id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// อัพเดท User
	user.Username = updatedUser.Username
	user.Role = updatedUser.Role
	r = db.Save(&user)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
