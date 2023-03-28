package main

import (
	"CRUD-API/handlers/user"
	// . "CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dsn = "postgres://navjsbdt:CXbvdzgydzdeZKUi_WYzMxzxAjJqnYbF@satao.db.elephantsql.com/navjsbdt"
var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

func main() {
	if err != nil {
		panic("failed to connect database")
	}
	r := gin.New()

	userHandler := user.NewUserHandler(db)
	r.GET("/user", userHandler.ListUser)
	r.POST("/user", userHandler.CreateUserHandler)
	r.DELETE("/user/:id", userHandler.DeleteUserHandler)

	r.Run()
}
