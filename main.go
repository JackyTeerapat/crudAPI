package main

import (
	"CRUD-API/handlers/assessment"
	"CRUD-API/handlers/assessment_article"
	"CRUD-API/handlers/assessment_progress"
	"CRUD-API/handlers/assessment_project"
	"CRUD-API/handlers/assessment_report"
	"CRUD-API/handlers/auth"
	"CRUD-API/handlers/degree"
	"CRUD-API/handlers/experience"
	"CRUD-API/handlers/exploration"
	"CRUD-API/handlers/position"
	"CRUD-API/handlers/profile"
	"CRUD-API/handlers/profile_attach"
	"CRUD-API/handlers/program"
	"CRUD-API/handlers/researcher"
	"CRUD-API/handlers/user"
	"CRUD-API/initializers"
	"CRUD-API/middlewares"

	// . "CRUD-API/models"

	_ "CRUD-API/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	initializers.LoadEnvVariables()
	db = initializers.ConnectDb()
}

// @title Researcher Service API
// @version 1.0.0
// @description This is a sample server for a researcher service.
// @host localhost:9000
// @BasePath /api/v1
func main() {
	r := gin.New()
	r.Use(middlewares.CORSMiddleware())
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := auth.NewAuthHandler(db)
	r.POST("/api/v1/signup", auth.SignUp)
	r.POST("/api/v1/login", auth.Login)

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

	//Article Zones
	assessmentHandler := assessment.NewAssessmentHandler(db)
	r.GET("/assessment/", assessmentHandler.ListAssessment)
	r.GET("/assessment/:id", assessmentHandler.GetAssessmentHandler)
	r.POST("/assessment", assessmentHandler.CreateAssessmentHandler)
	r.PUT("/assessment/:id", assessmentHandler.UpdateAssessmentHandler)
	r.DELETE("/assessment/:id", assessmentHandler.DeleteAssessmentHandler)

	//Article Zones
	articleHandler := assessment_article.NewArticleHandler(db)
	r.GET("/article/", articleHandler.ListArticle)
	r.GET("/article/:id", articleHandler.GetArticleHandler)
	r.POST("/article", articleHandler.CreateArticleHandler)
	r.PUT("/article/:id", articleHandler.UpdateArticleHandler)
	r.DELETE("/article/:id", articleHandler.DeleteArticleHandler)

	//AssessmentProgress Zones
	progressHandler := assessment_progress.NewProgressHandler(db)
	r.GET("progress/", progressHandler.ListProgress)
	r.GET("/progress/:id", progressHandler.GetProgressHandler)
	r.POST("/progress", progressHandler.CreateProgressHandler)
	r.PUT("/progress/:id", progressHandler.UpdateProgressHandler)
	r.DELETE("/progress/:id", progressHandler.DeleteProgressHandler)

	//AssessmentReport Zones
	reportHandler := assessment_report.NewReportHandler(db)
	r.GET("report/", reportHandler.ListReport)
	r.GET("/report/:id", reportHandler.GetReportHandler)
	r.POST("/report", reportHandler.CreateReportHandler)
	r.PUT("/report/:id", reportHandler.UpdateReportHandler)
	r.DELETE("/report/:id", reportHandler.DeleteReportHandler)

	//Assessment Project Zones
	projectHandler := assessment_project.NewProjectHandler(db)
	r.GET("/project", projectHandler.ListProjects)
	r.GET("/project/:id", projectHandler.GetProjectHandler)
	r.POST("/project", projectHandler.CreateProjectHandler)
	r.PUT("/project/:id", projectHandler.UpdateProjectHandler)
	r.DELETE("/project/:id", projectHandler.DeleteProjectHandler)

	//researcher
	researcherHandler := researcher.NewResearcherHandler(db)
	r.GET("/api/v1/researcher/profile_detail/:id", researcherHandler.ListResearcher)
	r.POST("/api/v1/researcher/profile", researcherHandler.CreateResearcher)
	r.PUT("/api/v1/researcher/profile/:id", researcherHandler.UpdateResearcher)
	r.DELETE("/api/v1/researcher/profile/:id", researcherHandler.DeleteResearcher)

	r.Run()

}
