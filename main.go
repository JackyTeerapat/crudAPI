package main

import (
	"CRUD-API/handlers/degree"
	"CRUD-API/handlers/experience"
	"CRUD-API/handlers/exploration"
	"CRUD-API/handlers/position"
	"CRUD-API/handlers/profile"
	"CRUD-API/handlers/profile_attach"
	"CRUD-API/handlers/program"
	"CRUD-API/handlers/researcher"
	"CRUD-API/handlers/user"

	// . "CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dsn = "postgres://navjsbdt:IrWX1ZnQiuYaMiTXCOwNCB-acHRKJgvT@satao.db.elephantsql.com/navjsbdt"
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

	//Experience Zones
	experienceHandler := experience.NewExperienceHandler(db)
	r.GET("experience/", experienceHandler.ListExperience)
	r.GET("/experience/:id", experienceHandler.GetExperienceHandler)
	r.POST("/experience", experienceHandler.CreateExperienceHandler)
	r.PUT("/experience/:id", experienceHandler.UpdateExperienceHandler)
	r.DELETE("/experience/:id", experienceHandler.DeleteExperienceHandler)

	//Exploration Zones
	explorationHandler := exploration.NewExplorationHandler(db)
	r.GET("exploration/", explorationHandler.ListExploration)
	r.GET("/exploration/:id", explorationHandler.GetExplorationHandler)
	r.POST("/exploration", explorationHandler.CreateExplorationHandler)
	r.PUT("/exploration/:id", explorationHandler.UpdateExplorationHandler)
	r.DELETE("/exploration/:id", explorationHandler.DeleteExplorationHandler)

	//Profile_attach Zones
	profile_attachHandler := profile_attach.NewProfile_attachHandler(db)
	r.GET("profile_attach/", profile_attachHandler.ListProfile_attach)
	r.GET("/profile_attach/:id", profile_attachHandler.GetProfile_attachHandler)
	r.POST("/profile_attach", profile_attachHandler.CreateProfile_attachHandler)
	r.PUT("/profile_attach/:id", profile_attachHandler.UpdateProfile_attachHandler)
	r.DELETE("/profile_attach/:id", profile_attachHandler.DeleteProfile_attachHandler)

	//researcher
	researcherHandler := researcher.NewResearcherHandler(db)
	r.GET("researcher/profile_detail/:id", researcherHandler.ListResearcher)
	r.Run()

}
