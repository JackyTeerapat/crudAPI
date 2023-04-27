package assessment

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	db *gorm.DB
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{db: db}
}
func (u *ProjectHandler) ListProjects(c *gin.Context) {
	var project []models.AssessmentProject

	r := u.db.Table("assessment_project").Find(&project)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

// GetProject godoc
// @Summary Get a project
// @Description Get a data user from database.
// @Tags Project
// @Produce  application/json
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project{}
// @Router /project/{id} [get]
func (u *ProjectHandler) GetProjectHandler(c *gin.Context) {
	var project models.AssessmentProject
	id := c.Param("id")
	r := u.db.Table("assessment_project").Where("id = ?", id).First(&project)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

// CreateProject godoc
// @Summary Create a project
// @Description Create a data project to database.
// @Tags Project
// @Produce  application/json
// @Param project body models.Project true "Project"
// @Success 200 {object} models.Project{}
// @Router /project [post]
func (u *ProjectHandler) CreateProjectHandler(c *gin.Context) {
	tx := u.db.Begin()
	var project models.AssessmentProject
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := tx.Create(&project)
	if err := r.Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a data project from database.
// @Tags Project
// @Produce  application/json
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project{}
// @Router /project/{id} [delete]
func (u *ProjectHandler) DeleteProjectHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง assessment project
		if err := u.db.Exec("TRUNCATE TABLE assessment_project CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER TABLE assessment_project ALTER COLUMN id SET DEFAULT nextval('assessment_project_id_seq'::regclass)").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All project data have been deleted."})
		return
	}

	// ลบข้อมูลตาม id ที่ระบุ
	r := u.db.Delete(&models.AssessmentProject{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("project with id %s has been deleted.", id)})
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update a data project to database.
// @Tags Project
// @Produce  application/json
// @Param id path int true "Project ID"
// @Param project body models.Project true "Project"
// @Success 200 {object} models.Project{}
// @Router /project/{id} [put]
func (u *ProjectHandler) UpdateProjectHandler(c *gin.Context) {
	var project models.Degree
	id := c.Param("id")

	//ตรวจสอบว่ามี degree นี้อยู่หรือไม่
	r := u.db.Table("assessment_project").Where("id = ?", id).First(&project)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล degree ด้วย ID ที่กำหนด
	r = u.db.Table("assessment_project").Where("id = ?", id).Updates(&project)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
