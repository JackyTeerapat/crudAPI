package main

import (
	"CRUD-API/handlers/position"
	"CRUD-API/handlers/user"
	"CRUD-API/handlers/degree"
	"CRUD-API/handlers/profile"
	"CRUD-API/handlers/program"
	"CRUD-API/handlers/assessment_progress"
	"CRUD-API/handlers/assessment_report"
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

	//User Zones
	userHandler := user.NewUserHandler(db)
	r.GET("/user", userHandler.ListUser)
	r.GET("/user/:id", userHandler.GetUserHandler)
	r.POST("/user", userHandler.CreateUserHandler)
	r.PUT("/user/:id", userHandler.UpdateUserHandler)
	r.DELETE("/user/:id", userHandler.DeleteUserHandler)

	//Position Zones
	positionHandler := position.NewPositionHandler(db)
	r.GET("/position", positionHandler.ListPosition)
	r.GET("/position/:id", positionHandler.GetPositionHandler)
	r.POST("/position", positionHandler.CreatePositionHandler)
	r.PUT("/position/:id", positionHandler.UpdatePositionHandler)
	r.DELETE("/position/:id", positionHandler.DeletePositionHandler)

	//Degree Zones
	degreeHandler := degree.NewDegreeHandler(db)
	r.GET("degree/", degreeHandler.ListDegree)
	r.GET("/degree/:id", degreeHandler.GetDegreeHandler)
	r.POST("/degree", degreeHandler.CreateDegreeHandler)
	r.PUT("/degree/:id", degreeHandler.UpdateDegreeHandler)
	r.DELETE("/degree/:id", degreeHandler.DeleteDegreeHandler)

	//Profile Zones
	profileHandler := profile.NewProfileHandler(db)
	r.GET("profile/", profileHandler.ListProfile)
	r.GET("/profile/:id", profileHandler.GetProfileHandler)
	r.POST("/profile", profileHandler.CreateProfileHandler)
	r.PUT("/profile/:id", profileHandler.UpdateProfileHandler)
	r.DELETE("/profile/:id", profileHandler.DeleteProfileHandler)

	
	//Program Zones
	programHandler := program.NewProgramHandler(db)
	r.GET("program/", programHandler.ListProgram)
	r.GET("/program/:id", programHandler.GetProgramHandler)
	r.POST("/program", programHandler.CreateProgramHandler)
	r.PUT("/program/:id", programHandler.UpdateProgramHandler)
	r.DELETE("/program/:id", programHandler.DeleteProgramHandler)


	//AssessmentProgress Zones
	progressHandler := assessment_progress.NewProgressHandler(db)
	r.GET("progress/", progressHandler.ListProgress)
	r.GET("/progress/:id", progressHandler.GetProgressHandler)
	r.POST("/progress", progressHandler.CreateProgressHandler)
	r.PUT("/progress/:id", progressHandler.UpdateProgressHandler)
	r.DELETE("/progress/:id", progressHandler.DeleteProgressHandler)

	//AssessmentReport Zones
	reportHandler := assessment_report.NewReportHandler(db)
	r.GET("progress/", reportHandler.ListReport)
	r.GET("/progress/:id", reportHandler.GetReportHandler)
	r.POST("/progress", reportHandler.CreateReportHandler)
	r.PUT("/progress/:id", reportHandler.UpdateReportHandler)
	r.DELETE("/progress/:id", reportHandler.DeleteReportHandler)



	r.Run()
}
