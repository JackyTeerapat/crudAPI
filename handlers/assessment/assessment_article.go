package assessment

import (
	"fmt"
	"net/http"

	"CRUD-API/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	db *gorm.DB
}

func NewArticleHandler(db *gorm.DB) *ArticleHandler {
	return &ArticleHandler{db: db}
}
func (u *ArticleHandler) ListArticle(c *gin.Context) {
	var articles []models.AssessmentArticle

	r := u.db.Table("assessment_article").Find(&articles)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, articles)
}
func (u *ArticleHandler) GetArticleHandler(c *gin.Context) {
	var article models.AssessmentArticle
	id := c.Param("id")
	r := u.db.Table("assessment_article").Where("id = ?", id).First(&article)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}

func (u *ArticleHandler) CreateArticleHandler(c *gin.Context) {
	var article models.AssessmentArticle

	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//รับมาแล้วสร้างเป็น ข้อมูล ลง Table
	r := u.db.Create(&article)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Status": "Success"})
}

func (u *ArticleHandler) DeleteArticleHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "all" {
		// ลบข้อมูลทั้งหมดในตาราง articles
		if err := u.db.Exec("TRUNCATE TABLE assessment_article CASCADE").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ตั้งค่า auto increment primary key เป็น 1
		if err := u.db.Exec("ALTER TABLE assessment_article ALTER COLUMN id SET DEFAULT nextval('assessment_article_id_seq'::regclass)").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All article data have been deleted."})
		return
	}

	// ลบข้อมูล article ตาม id ที่ระบุ
	r := u.db.Delete(&models.AssessmentArticle{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("article with id %s has been deleted.", id)})
}

func (u *ArticleHandler) UpdateArticleHandler(c *gin.Context) {
	var article models.AssessmentArticle
	id := c.Param("id")

	//ตรวจสอบว่ามี article นี้อยู่หรือไม่
	r := u.db.Table("assessment_article").Where("id = ?", id).First(&article)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	//แปลงข้อมูลที่ส่งเข้ามาเป็น JSON
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//อัปเดตข้อมูล article ด้วย ID ที่กำหนด
	r = u.db.Table("assessment_article").Where("id = ?", id).Updates(&article)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Success"})
}
