package position

import (
	"fmt"
	"net/http"

	"CRUD-API/api"
	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PositionHandler struct {
	db *gorm.DB
}

func NewPositionHandler(db *gorm.DB) *PositionHandler {
	return &PositionHandler{db: db}
}
func (u *PositionHandler) ListPosition(c *gin.Context) {
	var positions []models.Position

	r := u.db.Table("position").Find(&positions)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []map[string]interface{}
	for _, position := range positions {
		m := make(map[string]interface{})
		m["position_id"] = position.ID
		m["position_name"] = position.Position_name
		data = append(data, m)
	}

	// res := gin.H{
	// 	"description":  "SUCCESS",
	// 	"errorMessage": nil,
	// 	"status":       http.StatusOK,
	// 	"data":         data,
	// }

	res := api.ResponseApiWithDescription(http.StatusOK, data, "SUCCESS", nil)
	c.JSON(http.StatusOK, res)
}

// res := api.ResponseApiWithDescription(http.StatusOK, positions, "SUCCESS", nil)
// c.JSON(http.StatusOK, res)
// }

func (u *PositionHandler) GetPositionHandler(c *gin.Context) {
	var position models.Position
	id := c.Param("id")
	r := u.db.Table("position").Where("id = ?", id).First(&position)
	if r.RowsAffected == 0 {
		r := u.db.Create(&position)
		if err := r.Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error to create position": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"Status": "Success create position"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, position)
}

func (u *PositionHandler) CreatePositionHandler(c *gin.Context) {
	var position models.Position

	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&position)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *PositionHandler) DeletePositionHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง positions
		if err := u.db.Exec("TRUNCATE TABLE position CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER SEQUENCE position_id_seq RESTART WITH 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All position data have been deleted."})
		return
	}

	// ลบข้อมูล position ตาม id ที่ระบุ
	r := u.db.Delete(&models.Position{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("position with id %s has been deleted.", id)})
}

func (u *PositionHandler) UpdatePositionHandler(c *gin.Context) {
	var position models.Position
	id := c.Param("id")

	//ตรวจสอบว่ามี position นี้อยู่หรือไม่
	r := u.db.Table("position").Where("id = ?", id).First(&position)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล position ด้วย ID ที่กำหนด
	r = u.db.Table("position").Where("id = ?", id).Updates(&position)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
