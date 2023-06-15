package researcher_list

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"CRUD-API/api"
	"CRUD-API/models"
	"CRUD-API/models/custom"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	ResearcherList struct {
		db *gorm.DB
	}

	RequestInput struct {
		ResearcherName string `json:"researcher_name"`
		ProgramName    string `json:"program_name"`
		Page           int    `json:"page"`
		Limit          int    `json:"limit"`
	}

	ResearcherOutput struct {
		ProfileId     int                `json:"profile_id"`
		ProfileStatus string             `json:"profile_status"`
		FirstName     string             `json:"first_name"`
		LastName      string             `json:"last_name"`
		PositionId    int                `json:"position_id"`
		PositionName  string             `json:"position_name"`
		University    string             `json:"university"`
		Email         string             `json:"email"`
		PhoneNumber   string             `json:"phone_number"`
		Degree        []*ResponseDegree  `json:"degree"`
		Program       []*ResponseProgram `json:"program"`
	}

	ResponseDegree struct {
		Id               int    `json:"id" gorm:"column:id"`
		DegreeType       string `json:"degree_type" gorm:"column:degree_type"`
		DegreeProgram    string `json:"degree_program" gorm:"column:degree_program"`
		DegreeUniversity string `json:"degree_university" gorm:"column:degree_university"`
	}

	ResponseProgram struct {
		ProgramId   int    `json:"program_id" gorm:"column:id"`
		ProgramName string `json:"program_name" gorm:"column:program_name"`
	}

	ResponseDataContent struct {
		Content      []ResearcherOutput `json:"content"`
		TotalPage    int                `json:"total_page"`
		TotalObject  int                `json:"total_object"`
		CurrentPage  int                `json:"current_page"`
		IsLast       bool               `json:"is_last"`
		ObjectInPage int                `json:"object_in_page"`
	}
)

const (
	SELECT    = "SELECT "
	FROM      = " FROM "
	AS        = " AS "
	WHERE     = " WHERE "
	AND       = " AND "
	INNERJOIN = " INNER JOIN "
	UNION     = " UNION "
	ON        = " ON "
	ORDERBY   = " ORDER BY "
	GROUPBY   = " GROUP BY "
	ASC       = " ASC"
	DESC      = " DESC"
	NOTIN     = " NOT IN "
	OR        = " OR "
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

	page := req.Page - 1
	if page < 0 {
		page = 0
	}
	limit := 5
	if req.Limit != 0 {
		limit = req.Limit
	}
	page = page * limit
	sql := ""
	if req.ProgramName != "" {
		sql = generateQueryByProgramName(req)
	} else {
		sql = generateQuery(req)
	}

	total_count, err := CountTotalItem(sql, u, c)
	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	sql += " ORDER BY profile.id DESC OFFSET " + strconv.Itoa(page) + " ROWS FETCH NEXT " + strconv.Itoa(limit) + " ROWS ONLY"

	withContext := u.db.WithContext(c)

	var result []*custom.Profile
	tx := withContext.Raw(sql)
	err = tx.Find(&result).Error

	if err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	ids := make([]int, len(result))
	for i, v := range result {
		ids[i] = v.ProfileId
	}

	allDegree, err := findDegreeByProfileIds(ids, u, c)
	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	degreeMap := make(map[int][]*ResponseDegree)
	for _, v := range allDegree {
		degreeMap[v.Profile_id] = append(degreeMap[v.Profile_id], &ResponseDegree{
			Id:               v.ID,
			DegreeType:       v.Degree_type,
			DegreeProgram:    v.Degree_program,
			DegreeUniversity: v.Degree_university,
		})
	}

	allProgram, err := findProgramByProfileIds(ids, u, c)
	if err != nil {
		res := api.ResponseApi(http.StatusInternalServerError, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	programMap := make(map[int][]*ResponseProgram)
	for _, v := range allProgram {
		programMap[v.Profile_id] = append(programMap[v.Profile_id], &ResponseProgram{
			ProgramId:   v.Id,
			ProgramName: v.Program_name,
		})
	}

	var resDataContent ResponseDataContent
	count := 0
	for _, v := range result {
		resDataContent.Content = append(resDataContent.Content, ResearcherOutput{
			ProfileId:     v.ProfileId,
			ProfileStatus: v.ProfileStatus,
			FirstName:     v.FirstName,
			LastName:      v.LastName,
			PositionId:    v.PositionId,
			PositionName:  v.PositionName,
			University:    v.University,
			Email:         v.Email,
			PhoneNumber:   v.Phone_number,
			Degree:        degreeMap[v.ProfileId],
			Program:       programMap[v.ProfileId],
		})
		count++
	}
	resDataContent.IsLast = (page + limit) >= total_count
	resDataContent.CurrentPage = req.Page
	if req.Page < 1 {
		resDataContent.CurrentPage = req.Page + 1
	}
	resDataContent.TotalObject = total_count
	resDataContent.ObjectInPage = count
	resDataContent.TotalPage = total_count / limit
	if total_count%limit > 0 {
		resDataContent.TotalPage++
	}
	
	res := api.ResponseApi(http.StatusOK, resDataContent, nil)
	c.JSON(http.StatusOK, res)
}

func CountTotalItem(sqlStatement string, u *ResearcherList, ctx context.Context) (int, error) {
	count := 0
	sql := "SELECT COUNT(*) FROM (" + sqlStatement + ") T"
	withContext := u.db.WithContext(ctx)

	tx := withContext.Raw(sql)
	if err := tx.Find(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func generateQuery(req RequestInput) string {
	sql := &strings.Builder{}

	sql.WriteString(SELECT)
	sql.WriteString(" profile.id as profile_id,profile.profile_status,profile.first_name,profile.last_name,profile.university,profile.email,profile.phone_number, ")
	sql.WriteString(" position.id as position_id,position.position_name ")
	sql.WriteString(FROM)
	sql.WriteString("profile")
	sql.WriteString(INNERJOIN)

	sql.WriteString("position")
	sql.WriteString(ON)
	sql.WriteString("position.id = profile.position_id")
	sql.WriteString(WHERE)
	sql.WriteString("profile.profile_status")

	if req.ResearcherName != "" {
		lower := strings.ToLower(req.ResearcherName)
		sql.WriteString(AND)
		sql.WriteString("(LOWER(profile.first_name) = '" + lower + "' OR LOWER(profile.last_name) = '" + lower + "')")
	}
	return sql.String()
}

func generateQueryByProgramName(req RequestInput) string {
	sql := &strings.Builder{}

	sql.WriteString(SELECT)
	sql.WriteString(" profile.id as profile_id,profile.profile_status,profile.first_name,profile.last_name,profile.university,profile.email,profile.phone_number, ")
	sql.WriteString(" position.id as position_id,position.position_name ")
	sql.WriteString(FROM)
	sql.WriteString("profile")

	sql.WriteString(INNERJOIN)
	sql.WriteString("position")
	sql.WriteString(ON)
	sql.WriteString("position.id = profile.position_id")

	sql.WriteString(INNERJOIN)
	sql.WriteString("program")
	sql.WriteString(ON)
	sql.WriteString("program.profile_id = profile.id")

	sql.WriteString(WHERE)
	sql.WriteString("profile.profile_status")

	if req.ProgramName != "" {
		lower := strings.ToLower(req.ProgramName)
		sql.WriteString(AND)
		sql.WriteString("LOWER(program.program_name) = '" + lower + "'")
	}

	if req.ResearcherName != "" {
		lower := strings.ToLower(req.ResearcherName)
		sql.WriteString(AND)
		sql.WriteString("(LOWER(profile.first_name) = '" + lower + "' OR LOWER(profile.last_name) = '" + lower + "')")
	}
	return sql.String()
}

func findDegreeByProfileIds(ids []int, u *ResearcherList, ctx context.Context) (result []*models.Degree, err error) {
	condition := u.db.WithContext(ctx).Model(&models.Degree{}).Where("profile_id IN (?)", ids)
	err = condition.Find(&result).Error
	return result, err
}

func findProgramByProfileIds(ids []int, u *ResearcherList, ctx context.Context) (result []*models.Program, err error) {
	condition := u.db.WithContext(ctx).Model(&models.Program{}).Where("profile_id IN (?)", ids)
	err = condition.Find(&result).Error
	return result, err
}
