package user

import (
	"net/http"
	"strconv"
	"strings"

	"CRUD-API/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	UserList struct {
		db *gorm.DB
	}

	RequestInput struct {
		UserName string `json:"username"`
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
	}

	UserOutput struct {
		UserId   int    `json:"user_id"`
		UserName string `json:"username"`
		Role     string `json:"role"`
	}

	ResponseDataContent struct {
		Content     []UserOutput `json:"content"`
		TotalPage   int          `json:"total_page"`
		TotalObject int          `json:"total_object"`
		CurrentPage int          `json:"current_page"`
		IsLast      bool         `json:"is_last"`
	}
)

func UserListConnection(db *gorm.DB) *UserList {
	return &UserList{db: db}
}
func (u *UserList) ListUser(c *gin.Context) {
	var req RequestInput

	if err := c.BindJSON(&req); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	page := req.Page - 1
	if page < 0 {
		page = 0
	}
	limit := 5
	if req.Limit != 0 {
		limit = req.Limit
	}
	page = page * limit
	sqlQueryStatement := "ID, username, role"
	sqlStatement := "SELECT #STATEMENT# FROM users "

	if req.UserName != "" {
		sqlStatement += " WHERE username LIKE '%" + req.UserName + "%'"
	}

	total_count, err := CountTotalItem(sqlStatement, u)
	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	sqlStatement += " ORDER BY ID ASC OFFSET " + strconv.Itoa(page) + " ROWS FETCH NEXT " + strconv.Itoa(limit) + " ROWS ONLY"

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
		tmp := UserOutput{0, "", ""}
		if err := list.Scan(&tmp.UserId, &tmp.UserName, &tmp.Role); err != nil {
			res := api.ResponseApi(http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, res)
			return
		}
		resDataContent.Content = append(resDataContent.Content, tmp)
		count++
	}
	resDataContent.IsLast = (page + limit) >= count
	resDataContent.CurrentPage = req.Page
	if req.Page < 1 {
		resDataContent.CurrentPage = req.Page + 1
	}
	resDataContent.TotalObject = count
	resDataContent.TotalPage = total_count / limit
	if total_count%limit > 0 {
		resDataContent.TotalPage++
	}
	res := api.ResponseApi(http.StatusOK, resDataContent, nil)
	c.JSON(http.StatusOK, res)
}

func CountTotalItem(sqlStatement string, u *UserList) (int, error) {
	count := 0
	sqlStatement = strings.Replace(sqlStatement, "#STATEMENT#", "COUNT(*)", 1)
	row := u.db.Raw(sqlStatement).Row()
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
