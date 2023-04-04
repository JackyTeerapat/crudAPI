package researcher_list

import (
	"net/http"
	"strconv"
	"strings"

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
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
}

type ResearcherOutput struct {
	ResearcherName string `json:"researcher_name"`
	University     string `json:"university"`
	ExploreYear    string `json:"explore_year"`
	ProjectTitle   string `json:"project_title"`
	ResearcherId   int    `json:"researcher_id"`
}

type Response struct {
	Status       int                 `json:"status"`
	Description  string              `json:"description"`
	ErrorMessage string              `json:"errorMessage"`
	Data         ResponseDataContent `json:"data"`
}

type ResponseDataContent struct {
	Content     []ResearcherOutput `json:"content"`
	TotalPage   int                `json:"total_page"`
	TotalObject int                `json:"total_object"`
	CurrentPage int                `json:"current_page"`
	IsLast      bool               `json:"is_last"`
}

func ResearcherListConnection(db *gorm.DB) *ResearcherList {
	return &ResearcherList{db: db}
}
func (u *ResearcherList) ListResearcher(c *gin.Context) {
	var req RequestInput

	var res Response
	if err := c.BindJSON(&req); err != nil {
		res.Status = http.StatusBadRequest
		res.ErrorMessage = err.Error()
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
	sqlCountStatement := "COUNT(*)"
	sqlQueryStatement := "profile.id as researcher_id, profile.first_name, profile.last_name, profile.university, exploration.explore_year, assessment_project.project_title"
	sqlStatement := "SELECT #STATEMENT# FROM profile " +
		"JOIN exploration ON profile.id = exploration.profile_id " +
		"JOIN assessment ON profile.id = assessment.profile_id " +
		"JOIN assessment_project ON assessment.project_id = assessment_project.id "

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
		if isAddWhere {
			sqlStatement += " AND exploration.explore_year LIKE '%" + req.ExploreYear + "%'"
		} else {
			isAddWhere = true
			sqlStatement += " WHERE exploration.explore_year LIKE '%" + req.ExploreYear + "%'"
		}
	}

	if req.ProjectTitle != "" {
		if isAddWhere {
			sqlStatement += " AND assessment_project.project_title LIKE '%" + req.ProjectTitle + "%'"
		} else {
			isAddWhere = true
			sqlStatement += " WHERE assessment_project.project_title LIKE '%" + req.ProjectTitle + "%'"
		}
	}

	count := CountTotalItem(sqlStatement, sqlCountStatement, u)
	sqlStatement += " ORDER BY profile.id DESC OFFSET " + strconv.Itoa(page) + " ROWS FETCH NEXT " + strconv.Itoa(limit) + " ROWS ONLY"

	sqlStatement = strings.Replace(sqlStatement, "#STATEMENT#", sqlQueryStatement, 1)
	list, err := u.db.Raw(sqlStatement).Rows()
	defer list.Close()

	if err != nil {
		res.Status = http.StatusBadRequest
		res.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, res)
		//c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("An error occurred while fetching data %v", err.Error())})
		return
	}

	var resDataContent ResponseDataContent
	for list.Next() {
		tmp := ResearcherOutput{"", "", "", "", 0}
		first, last := "", ""
		if err := list.Scan(&tmp.ResearcherId, &first, &last, &tmp.University, &tmp.ExploreYear, &tmp.ProjectTitle); err != nil {
			res.Status = http.StatusBadRequest
			res.ErrorMessage = err.Error()
			c.JSON(http.StatusBadRequest, res)
			return
		}
		tmp.ResearcherName = first + " " + last
		resDataContent.Content = append(resDataContent.Content, tmp)
	}
	res.Status = http.StatusOK
	res.Description = "SUCCESS"
	resDataContent.IsLast = (page + limit) >= count
	resDataContent.CurrentPage = req.Page
	resDataContent.TotalObject = count
	resDataContent.TotalPage = count / limit
	if count%limit > 0 {
		resDataContent.TotalPage++
	}
	res.Data = resDataContent
	c.JSON(http.StatusOK, res)
}

func CountTotalItem(sqlStatement, sqlCountStatement string, u *ResearcherList) int {
	count := 0
	sqlStatement = strings.Replace(sqlStatement, "#STATEMENT#", sqlCountStatement, 1)
	row := u.db.Raw(sqlStatement).Row()
	err := row.Scan(&count)
	if err != nil {
		panic(err)
	}

	return count
}
