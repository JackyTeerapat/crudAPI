package researcher_list

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResearcherList struct {
	db *gorm.DB
}

type RequestInput struct {
	ResearcherName string `json:"researcher_name"`
	University     string `json:"university"`
	ExploreYear    string `json:"explore_year"`
	ProjectTitle   string `json:"project_title"`
	Page           string `json:"page"`
	Limit          string `json:"limit"`
}

func ResearcherListConnection(db *gorm.DB) *ResearcherList {
	return &ResearcherList{db: db}
}
func (u *ResearcherList) ListPosition(c *gin.Context) {
	var req RequestInput

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "test")
}
