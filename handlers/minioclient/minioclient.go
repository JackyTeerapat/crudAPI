package minioclient

import (
	"CRUD-API/api"
	"CRUD-API/models"
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

type MinioClient struct {
	mc         *minio.Client
	bucketName string
	db         *gorm.DB
}

type FileResponseData struct {
	FileName        string `json:"file_name"`
	FileType        string `json:"file_type"`
	FileUploadError string `json:"file_upload_error"`
}

type UploadedFile struct {
	FileName   string `json:"file_name"`
	FileBase64 string `json:"base64"`
}

func MinioClientConnect(db *gorm.DB) *MinioClient {
	endPoint := os.Getenv("MINIO_ENDPOINT")
	useSSL := true
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	minioClient, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &MinioClient{minioClient, bucketName, db}
}

func (m *MinioClient) UploadFile(c *gin.Context) {
	//multi file
	ctx := context.Background()
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	files := form.File["uploadfile"]

	directory := form.Value["directory_file"]
	directory_id := form.Value["directory_id"]

	resData := []FileResponseData{}
	for i, file := range files {

		if i >= len(directory_id) {
			continue
		}
		if i >= len(directory) {
			continue
		}
		row_id, r_err := strconv.Atoi(directory_id[i])
		if r_err != nil {
			continue
		}
		timenow := time.Now()
		timestamp := timenow.Format("20060102-15040506")
		fileName := timestamp + "-" + file.Filename
		newName := directory[i] + "/" + timestamp + "-" + file.Filename

		contentType := file.Header.Values("Content-Type")
		mimeType := "application/octet-stream"
		if len(contentType) > 0 {
			mimeType = contentType[0]
		}
		fileBuffer, err := file.Open()
		if err != nil {
			log.Panic(err)
		}

		fileBuffer.Close()

		fileData := FileResponseData{
			FileName: fileName,
			FileType: directory[i],
		}

		var u_err error
		switch directory[i] {
		case "assessment":
			_, u_err = UpdateAssessment(m.db, fileName, directory[i], row_id, false)
		case "project":
			_, u_err = UpdateAssessmentProject(m.db, fileName, directory[i], row_id, false)
		case "progress":
			_, u_err = UpdateAssessmentProgress(m.db, fileName, directory[i], row_id, false)
		case "report":
			_, u_err = UpdateAssessmentReport(m.db, fileName, directory[i], row_id, false)
		case "article":
			_, u_err = UpdateAssessmentArticle(m.db, fileName, directory[i], row_id, false)
		default:
			u_err = UpSertProfileAttach(m.db, fileName, directory[i], row_id)
		}
		if u_err != nil {
			fileData.FileUploadError = u_err.Error()
		} else {
			if _, err := m.mc.PutObject(ctx, m.bucketName, newName, fileBuffer, file.Size, minio.PutObjectOptions{ContentType: mimeType}); err != nil {
				fileData.FileUploadError = err.Error()
			}
		}

		resData = append(resData, fileData)
	}
	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusCreated, res)
	return
}

func (m *MinioClient) UploadFileBase64(c *gin.Context) {
	// single file
	ctx := context.Background()

	var req []UploadedFile

	if err := c.BindJSON(&req); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	directory := c.Param("directory")
	if directory == "" {
		directory = "upload"
	}

	resData := []FileResponseData{}
	for _, v := range req {

		timenow := time.Now()
		timestamp := timenow.Format("20060102-15040506")
		fileName := timestamp + "-" + v.FileName
		newName := directory + "/" + timestamp + "-" + v.FileName

		arr := strings.Split(v.FileBase64, ",")
		mimeType := "application/octet-stream"

		if matched, _ := regexp.MatchString(":(.*?);", arr[0]); matched {
			r, _ := regexp.Compile(":(.*?);")
			mimeType = r.FindStringSubmatch(arr[0])[1]
		}

		fileByte, err := base64.StdEncoding.DecodeString(arr[1])
		if err != nil {
			log.Panic(err)
		}
		fileBuffer := bytes.NewReader(fileByte)
		size := int64(len(fileByte))

		if _, err := m.mc.PutObject(ctx, m.bucketName, newName, fileBuffer, size, minio.PutObjectOptions{ContentType: mimeType}); err != nil {
			res := api.ResponseApi(http.StatusInternalServerError, nil, err)
			c.JSON(http.StatusInternalServerError, res)
			return
		} else {
			fileData := FileResponseData{
				FileName: fileName,
				FileType: directory,
			}
			resData = append(resData, fileData)
		}
	}
	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusCreated, res)
	return
}

func (m *MinioClient) DeleteFile(c *gin.Context) {
	ctx := context.Background()

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	directory := form.Value["directory_file"]
	directory_id := form.Value["directory_id"]

	resData := []FileResponseData{}
	for i, v := range directory {
		if i >= len(directory_id) {
			continue
		}
		row_id, r_err := strconv.Atoi(directory_id[i])
		if r_err != nil {
			continue
		}

		filename := ""
		var db_err error
		switch directory[i] {
		case "assessment":
			filename, db_err = UpdateAssessment(m.db, "", directory[i], row_id, true)
		case "project":
			filename, db_err = UpdateAssessmentProject(m.db, "", directory[i], row_id, true)
		case "progress":
			filename, db_err = UpdateAssessmentProgress(m.db, "", directory[i], row_id, true)
		case "report":
			filename, db_err = UpdateAssessmentReport(m.db, "", directory[i], row_id, true)
		case "article":
			filename, db_err = UpdateAssessmentArticle(m.db, "", directory[i], row_id, true)
		default:
			filename, db_err = DeleteProfileAttach(m.db, directory[i], row_id)
		}
		if db_err != nil {
			res := api.ResponseApi(http.StatusInternalServerError, nil, err)
			c.JSON(http.StatusInternalServerError, res)
			return
		} else {
			opts := minio.RemoveObjectOptions{GovernanceBypass: true}
			err := m.mc.RemoveObject(ctx, m.bucketName, v+"/"+filename, opts)
			if err != nil {
				res := api.ResponseApi(http.StatusInternalServerError, nil, err)
				c.JSON(http.StatusInternalServerError, res)
				return
			}
		}
		fileData := FileResponseData{
			FileName: filename,
			FileType: directory[i],
		}

		resData = append(resData, fileData)

	}

	//res := api.ResponseApi(http.StatusOK, "Success", nil)
	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusOK, res)
	return
}

func UpSertProfileAttach(db *gorm.DB, filename, data_type string, profile_id int) error {

	var profile_attach models.Profile_attach

	// check exist
	r := db.Table("profile_attach").Where("profile_id = ?", profile_id).Where("File_action = ?", data_type).First(&profile_attach)
	if r.RowsAffected == 0 {
		profile_attach.File_name = filename
		profile_attach.File_storage = data_type + "/" + filename
		profile_attach.File_action = data_type
		profile_attach.Activated = true
		profile_attach.Profile_id = profile_id
		r = db.Create(&profile_attach)
		if err := r.Error; err != nil {
			return err
		}
	} else {
		//update
		profile_attach.File_name = filename
		profile_attach.File_storage = data_type + "/" + filename
		profile_attach.Activated = true
		profile_attach.UpdatedAt = time.Now()
		r = db.Table("profile_attach").Updates(&profile_attach)
		if err := r.Error; err != nil {
			return err
		}
	}
	return nil
}

func UpdateAssessment(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var assessment models.Assessment
	res_flie_name := ""
	// check exist
	r := db.Table("assessment").Where("id = ?", profile_id).First(&assessment)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = assessment.Assessment_file_name
			assessment.Assessment_file_name = ""
			assessment.Assessment_file_storage = ""
		} else {
			assessment.Assessment_file_name = filename
			assessment.Assessment_file_storage = data_type + "/" + filename
		}
		assessment.UpdatedAt = time.Now()
		r = db.Table("assessment").Select("assessment_file_name", "assessment_file_storage", "updated_at").Updates(&assessment)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func UpdateAssessmentProject(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var project models.Project
	res_flie_name := ""
	// check exist
	r := db.Table("assessment_project").Where("id = ?", profile_id).First(&project)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = project.File_name
			project.File_name = ""
			project.File_storage = ""
		} else {
			project.File_name = filename
			project.File_storage = data_type + "/" + filename
		}
		project.Updated_at = time.Now()
		r = db.Table("assessment_project").Select("file_name", "file_storage", "updated_at").Updates(&project)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func UpdateAssessmentProgress(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var progress models.Progress
	res_flie_name := ""

	// check exist
	r := db.Table("assessment_progress").Where("id = ?", profile_id).First(&progress)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = progress.File_name
			progress.File_name = ""
			progress.File_storage = ""
		} else {
			progress.File_name = filename
			progress.File_storage = data_type + "/" + filename
		}
		progress.UpdatedAt = time.Now()
		r = db.Table("assessment_progress").Select("file_name", "file_storage", "updated_at").Updates(&progress)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func UpdateAssessmentReport(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var repport models.Report
	res_flie_name := ""

	// check exist
	r := db.Table("assessment_report").Where("id = ?", profile_id).First(&repport)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = repport.File_name
			repport.File_name = ""
			repport.File_storage = ""
		} else {
			repport.File_name = filename
			repport.File_storage = data_type + "/" + filename
		}
		repport.UpdatedAt = time.Now()
		r = db.Table("assessment_report").Select("file_name", "file_storage", "updated_at").Updates(&repport)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func UpdateAssessmentArticle(db *gorm.DB, filename, data_type string, profile_id int, is_delete bool) (string, error) {

	var article models.Article
	res_flie_name := ""

	// check exist
	r := db.Table("assessment_article").Where("id = ?", profile_id).First(&article)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		if is_delete {
			res_flie_name = article.File_name
			article.File_name = ""
			article.File_storage = ""
		} else {
			article.File_name = filename
			article.File_storage = data_type + "/" + filename
		}
		article.UpdatedAt = time.Now()
		r = db.Table("assessment_article").Select("file_name", "file_storage", "updated_at").Updates(&article)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}

func DeleteProfileAttach(db *gorm.DB, data_type string, profile_id int) (string, error) {

	var profile_attach models.Profile_attach
	res_flie_name := ""
	// check exist
	r := db.Table("profile_attach").Where("profile_id = ?", profile_id).Where("File_action = ?", data_type).First(&profile_attach)
	if r.RowsAffected == 0 {
		if err := r.Error; err != nil {
			return "", err
		}
	} else {
		//update
		res_flie_name = profile_attach.File_name
		profile_attach.File_name = ""
		profile_attach.File_storage = ""
		profile_attach.Activated = false
		profile_attach.UpdatedAt = time.Now()
		r = db.Table("profile_attach").Select("file_name", "file_storage", "activated", "updated_at").Updates(&profile_attach)
		if err := r.Error; err != nil {
			return "", err
		}
	}
	return res_flie_name, nil
}
