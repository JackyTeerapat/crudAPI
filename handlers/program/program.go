package program

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProgramHandler struct {
	db *gorm.DB
}

func NewProgramHandler(db *gorm.DB) *ProgramHandler {
	return &ProgramHandler{db: db}
}
func (u *ProgramHandler) ListProgram(c *gin.Context) {
	var programs []models.Program

	r := u.db.Table("program").Find(&programs)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, programs)
}
func (u *ProgramHandler) GetProgramHandler(c *gin.Context) {
	var program models.Program
	id := c.Param("id")
	r := u.db.Table("program").Where("id = ?", id).First(&program)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "program not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, program)
}

func (u *ProgramHandler) CreateProgramHandler(c *gin.Context) {
	var program models.Program

	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&program)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *ProgramHandler) DeleteProgramHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง programs
		if err := u.db.Exec("TRUNCATE TABLE program RESTART IDENTITY CASCADE").Error; err != nil {		
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER SEQUENCE program RESTART WITH 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All program data have been deleted."})
		return
	}

	// ลบข้อมูล program ตาม id ที่ระบุ
	r := u.db.Delete(&models.Program{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("program with id %s has been deleted.", id)})
}

func (u *ProgramHandler) UpdateProgramHandler(c *gin.Context) {
	var program models.Program
	id := c.Param("id")

	//ตรวจสอบว่ามี program นี้อยู่หรือไม่
	r := u.db.Table("program").Where("id = ?", id).First(&program)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "program not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล program ด้วย ID ที่กำหนด
	r = u.db.Table("program").Where("id = ?", id).Updates(&program)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}

