package researcher_list

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"CRUD-API/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	ResearcherList struct {
		db *gorm.DB
	}

	RequestInput struct {
		ResearcherName string `json:"researcher_name"`
		University     string `json:"university"`
		ExploreYear    string `json:"explore_year"`
		ProjectTitle   string `json:"project_title"`
		Page           int    `json:"page"`
		Limit          int    `json:"limit"`
	}

	ResearcherOutput struct {
		ResearcherName string `json:"researcher_name"`
		University     string `json:"university"`
		ExploreYear    string `json:"explore_year"`
		ProjectTitle   string `json:"project_title"`
		ResearcherId   int    `json:"researcher_id"`
	}

	ResponseDataContent struct {
		Content     []ResearcherOutput `json:"content"`
		TotalPage   int                `json:"total_page"`
		TotalObject int                `json:"total_object"`
		CurrentPage int                `json:"current_page"`
		IsLast      bool               `json:"is_last"`
	}
)

func ResearcherListConnection(db *gorm.DB) *ResearcherList {
	return &ResearcherList{db: db}
}
func (u *ResearcherList) ListResearcher(c *gin.Context) {
	var req RequestInput

	if err := c.BindJSON(&req); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	page := req.Page
	limit := 5
	if req.Limit != 0 {
		limit = req.Limit
	}
	page = page * limit
	isAddWhere := false
	sqlQueryStatement := "profile.id as researcher_id, profile.first_name, profile.last_name, profile.university, exp.explore_year, assessment_project.project_title"
	sqlStatement := "SELECT #STATEMENT# FROM (SELECT profile_id, MAX(explore_year) as explore_year FROM exploration #EXPWHERESTATEMENT# GROUP BY profile_id) exp " +
		"JOIN assessment ON assessment.profile_id = exp.profile_id " +
		"JOIN assessment_project ON assessment.project_id = assessment_project.id " +
		"JOIN profile ON profile.id = exp.profile_id "

	if req.ResearcherName != "" {
		isAddWhere = true
		sqlStatement += " WHERE (profile.first_name LIKE '%" + req.ResearcherName + "%' OR profile.last_name LIKE '%" + req.ResearcherName + "%')"
	}

	if req.University != "" {
		if isAddWhere {
			sqlStatement += " AND profile.university LIKE '%" + req.University + "%'"
		} else {
			isAddWhere = true
			sqlStatement += " WHERE profile.university LIKE '%" + req.University + "%'"
		}
	}

	if req.ExploreYear != "" {
		exp_where := " WHERE explore_year LIKE '%" + req.ExploreYear + "%'"
		sqlStatement = strings.Replace(sqlStatement, "#EXPWHERESTATEMENT#", exp_where, 1)
	} else {
		sqlStatement = strings.Replace(sqlStatement, "#EXPWHERESTATEMENT#", "", 1)
	}

	if req.ProjectTitle != "" {
		if isAddWhere {
			sqlStatement += " AND assessment_project.project_title LIKE '%" + req.ProjectTitle + "%'"
		} else {
			isAddWhere = true
			sqlStatement += " WHERE assessment_project.project_title LIKE '%" + req.ProjectTitle + "%'"
		}
	}

	total_count, err := CountTotalItem(sqlStatement, u)
	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	sqlStatement += " ORDER BY profile.id DESC OFFSET " + strconv.Itoa(page) + " ROWS FETCH NEXT " + strconv.Itoa(limit) + " ROWS ONLY"

	sqlStatement = strings.Replace(sqlStatement, "#STATEMENT#", sqlQueryStatement, 1)
	list, err := u.db.Raw(sqlStatement).Rows()
	defer list.Close()

	if err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		//c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("An error occurred while fetching data %v", err.Error())})
		return
	}

	var resDataContent ResponseDataContent
	count := 0
	for list.Next() {
		tmp := ResearcherOutput{"", "", "", "", 0}
		first, last := "", ""
		if err := list.Scan(&tmp.ResearcherId, &first, &last, &tmp.University, &tmp.ExploreYear, &tmp.ProjectTitle); err != nil {
			res := api.ResponseApi(http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, res)
			return
		}
		tmp.ResearcherName = first + " " + last
		resDataContent.Content = append(resDataContent.Content, tmp)
		count++
	}
	resDataContent.IsLast = (page + limit) >= count
	resDataContent.CurrentPage = req.Page
	resDataContent.TotalObject = count
	resDataContent.TotalPage = total_count / limit
	if count%limit > 0 {
		resDataContent.TotalPage++
	}
	res := api.ResponseApi(http.StatusOK, resDataContent, nil)
	c.JSON(http.StatusOK, res)
}

func CountTotalItem(sqlStatement string, u *ResearcherList) (int, error) {
	count := 0
	sqlStatement = strings.Replace(sqlStatement, "#STATEMENT#", "COUNT(*)", 1)
	fmt.Println(sqlStatement)
	row := u.db.Raw(sqlStatement).Row()
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
