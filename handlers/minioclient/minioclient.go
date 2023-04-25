package minioclient

import (
	"CRUD-API/api"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

type (
	MinioClient struct {
		mc         *minio.Client
		bucketName string
		db         *gorm.DB
	}

	FileResponseData struct {
		FileName string `json:"file_name"`
		FileType string `json:"file_type"`
		FileUrl  string `json:"file_url"`
	}

	MinioInput struct {
		FileName      string `json:"file_name"`
		FileBase64    string `json:"base64"`
		DirectoryFile string `json:"directory_file"`
		DirectoryId   int    `json:"directory_id"`
	}
)

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
	ctx := c
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	files := form.File["uploadfile"]

	profile_id := -1
	assessment_id := -1
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
		target := directory[i] + "/" + fileName

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
			assessment_id = row_id
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
			profile_id = row_id
			u_err = UpSertProfileAttach(m.db, fileName, directory[i], row_id)
		}
		if u_err != nil {
			res := api.ResponseApi(http.StatusBadRequest, nil, u_err)
			c.JSON(http.StatusBadRequest, res)
			if profile_id != -1 {
				DeleteProfile(m.db, profile_id)
			}
			if assessment_id != -1 {
				DeleteAssessment(m.db, assessment_id)
			}
			RollbackDeleteFile(c, m, resData)
			return
		} else {
			if _, err := m.mc.PutObject(ctx, m.bucketName, target, fileBuffer, file.Size, minio.PutObjectOptions{ContentType: mimeType}); err != nil {
				res := api.ResponseApi(http.StatusBadRequest, nil, err)
				c.JSON(http.StatusBadRequest, res)
				if profile_id != -1 {
					DeleteProfile(m.db, profile_id)
				}
				if assessment_id != -1 {
					DeleteAssessment(m.db, assessment_id)
				}
				RollbackDeleteFile(c, m, resData)
				return
			}
		}

		resData = append(resData, fileData)
	}
	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusOK, res)
	return
}

func (m *MinioClient) UploadFileBase64(c *gin.Context) {

	ctx := c

	var req []MinioInput
	profile_id := -1
	assessment_id := -1

	if err := c.BindJSON(&req); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	resData := []FileResponseData{}
	for _, v := range req {

		timenow := time.Now()
		timestamp := timenow.Format("20060102-15040506")
		fileName := timestamp + "-" + v.FileName
		target := v.DirectoryFile + "/" + fileName

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

		fileData := FileResponseData{
			FileName: fileName,
			FileType: v.DirectoryFile,
		}

		var u_err error
		switch v.DirectoryFile {
		case "assessment":
			assessment_id = v.DirectoryId
			_, u_err = UpdateAssessment(m.db, fileName, v.DirectoryFile, v.DirectoryId, false)
		case "project":
			_, u_err = UpdateAssessmentProject(m.db, fileName, v.DirectoryFile, v.DirectoryId, false)
		case "progress":
			_, u_err = UpdateAssessmentProgress(m.db, fileName, v.DirectoryFile, v.DirectoryId, false)
		case "report":
			_, u_err = UpdateAssessmentReport(m.db, fileName, v.DirectoryFile, v.DirectoryId, false)
		case "article":
			_, u_err = UpdateAssessmentArticle(m.db, fileName, v.DirectoryFile, v.DirectoryId, false)
		default:
			profile_id = v.DirectoryId
			u_err = UpSertProfileAttach(m.db, fileName, v.DirectoryFile, v.DirectoryId)
		}

		if u_err != nil {
			res := api.ResponseApi(http.StatusBadRequest, nil, u_err)
			c.JSON(http.StatusBadRequest, res)
			if profile_id != -1 {
				DeleteProfile(m.db, profile_id)
			}
			if assessment_id != -1 {
				DeleteAssessment(m.db, assessment_id)
			}
			RollbackDeleteFile(c, m, resData)
			return
		} else {
			if _, err := m.mc.PutObject(ctx, m.bucketName, target, fileBuffer, size, minio.PutObjectOptions{ContentType: mimeType}); err != nil {
				res := api.ResponseApi(http.StatusBadRequest, nil, err)
				c.JSON(http.StatusBadRequest, res)
				if profile_id != -1 {
					DeleteProfile(m.db, profile_id)
				}
				if assessment_id != -1 {
					DeleteAssessment(m.db, assessment_id)
				}
				RollbackDeleteFile(c, m, resData)
				return
			}
		}

		resData = append(resData, fileData)
	}
	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusOK, res)
	return
}

func (m *MinioClient) UploadUpdateFile(c *gin.Context) {
	//multi file
	ctx := c
	tx := m.db.Begin()

	form, err := c.MultipartForm()
	if err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
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
		target := directory[i] + "/" + fileName

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
			_, u_err = UpdateAssessment(tx, fileName, directory[i], row_id, false)
		case "project":
			_, u_err = UpdateAssessmentProject(tx, fileName, directory[i], row_id, false)
		case "progress":
			_, u_err = UpdateAssessmentProgress(tx, fileName, directory[i], row_id, false)
		case "report":
			_, u_err = UpdateAssessmentReport(tx, fileName, directory[i], row_id, false)
		case "article":
			_, u_err = UpdateAssessmentArticle(tx, fileName, directory[i], row_id, false)
		default:
			u_err = UpSertProfileAttach(tx, fileName, directory[i], row_id)
		}
		if u_err != nil {
			res := api.ResponseApi(http.StatusBadRequest, nil, u_err)
			c.JSON(http.StatusBadRequest, res)
			tx.Rollback()
			RollbackDeleteFile(c, m, resData)
			return
		} else {
			if _, err := m.mc.PutObject(ctx, m.bucketName, target, fileBuffer, file.Size, minio.PutObjectOptions{ContentType: mimeType}); err != nil {
				res := api.ResponseApi(http.StatusBadRequest, nil, err)
				c.JSON(http.StatusBadRequest, res)
				tx.Rollback()
				RollbackDeleteFile(c, m, resData)
				return
			}
		}

		resData = append(resData, fileData)
	}
	tx.Commit()
	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusOK, res)
	return
}

func (m *MinioClient) UploadUpdateFileBase64(c *gin.Context) {
	ctx := c

	tx := m.db.Begin()
	var req []MinioInput

	if err := c.BindJSON(&req); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	resData := []FileResponseData{}
	for _, v := range req {

		timenow := time.Now()
		timestamp := timenow.Format("20060102-15040506")
		fileName := timestamp + "-" + v.FileName
		target := v.DirectoryFile + "/" + fileName

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

		fileData := FileResponseData{
			FileName: fileName,
			FileType: v.DirectoryFile,
		}

		var u_err error
		switch v.DirectoryFile {
		case "assessment":
			_, u_err = UpdateAssessment(tx, fileName, v.DirectoryFile, v.DirectoryId, false)
		case "project":
			_, u_err = UpdateAssessmentProject(tx, fileName, v.DirectoryFile, v.DirectoryId, false)
		case "progress":
			_, u_err = UpdateAssessmentProgress(tx, fileName, v.DirectoryFile, v.DirectoryId, false)
		case "report":
			_, u_err = UpdateAssessmentReport(tx, fileName, v.DirectoryFile, v.DirectoryId, false)
		case "article":
			_, u_err = UpdateAssessmentArticle(tx, fileName, v.DirectoryFile, v.DirectoryId, false)
		default:
			u_err = UpSertProfileAttach(tx, fileName, v.DirectoryFile, v.DirectoryId)
		}

		if u_err != nil {
			res := api.ResponseApi(http.StatusBadRequest, nil, u_err)
			c.JSON(http.StatusBadRequest, res)
			tx.Rollback()
			RollbackDeleteFile(c, m, resData)
			return
		} else {
			if _, err := m.mc.PutObject(ctx, m.bucketName, target, fileBuffer, size, minio.PutObjectOptions{ContentType: mimeType}); err != nil {
				res := api.ResponseApi(http.StatusBadRequest, nil, err)
				c.JSON(http.StatusBadRequest, res)
				tx.Rollback()
				RollbackDeleteFile(c, m, resData)
				return
			}
		}

		resData = append(resData, fileData)
	}
	tx.Commit()
	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusOK, res)
	return
}

func (m *MinioClient) GetFile(c *gin.Context) {
	ctx := c

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	directory := form.Value["directory_file"]
	directory_id := form.Value["directory_id"]

	resData := []FileResponseData{}
	for i := range directory {
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
			filename, db_err = GetAssessment(m.db, directory[i], row_id)
		case "project":
			filename, db_err = GetAssessmentProject(m.db, directory[i], row_id)
		case "progress":
			filename, db_err = GetAssessmentProgress(m.db, directory[i], row_id)
		case "report":
			filename, db_err = GetAssessmentReport(m.db, directory[i], row_id)
		case "article":
			filename, db_err = GetAssessmentArticle(m.db, directory[i], row_id)
		default:
			filename, db_err = GetProfileAttach(m.db, directory[i], row_id)
		}

		resurl := ""
		if db_err != nil {
			res := api.ResponseApi(http.StatusInternalServerError, nil, err)
			c.JSON(http.StatusInternalServerError, res)
			return
		} else {
			reqParams := make(url.Values)
			presignedURL, err := m.mc.PresignedGetObject(ctx, m.bucketName, directory[i]+"/"+filename, time.Second*24*60*60, reqParams)
			if err != nil {
				res := api.ResponseApi(http.StatusInternalServerError, nil, err)
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			resurl = fmt.Sprintf("%s", presignedURL)
		}
		fileData := FileResponseData{
			FileName: filename,
			FileType: directory[i],
			FileUrl:  resurl,
		}

		resData = append(resData, fileData)

	}

	//res := api.ResponseApi(http.StatusOK, "Success", nil)
	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusOK, res)
	return
}

func (m *MinioClient) DeleteFile(c *gin.Context) {
	ctx := c
	var req []MinioInput
	if err := c.BindJSON(&req); err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	resData := []FileResponseData{}
	for i := range req {

		row_id := req[i].DirectoryId

		filename := ""
		var db_err error
		switch req[i].DirectoryFile {
		case "assessment":
			filename, db_err = UpdateAssessment(m.db, "", req[i].DirectoryFile, row_id, true)
		case "project":
			filename, db_err = UpdateAssessmentProject(m.db, "", req[i].DirectoryFile, row_id, true)
		case "progress":
			filename, db_err = UpdateAssessmentProgress(m.db, "", req[i].DirectoryFile, row_id, true)
		case "report":
			filename, db_err = UpdateAssessmentReport(m.db, "", req[i].DirectoryFile, row_id, true)
		case "article":
			filename, db_err = UpdateAssessmentArticle(m.db, "", req[i].DirectoryFile, row_id, true)
		default:
			filename, db_err = DeleteProfileAttach(m.db, req[i].DirectoryFile, row_id)
		}
		if db_err != nil {
			res := api.ResponseApi(http.StatusInternalServerError, nil, db_err)
			c.JSON(http.StatusInternalServerError, res)
			return
		} else {
			opts := minio.RemoveObjectOptions{GovernanceBypass: true}
			err := m.mc.RemoveObject(ctx, m.bucketName, req[i].DirectoryFile+"/"+filename, opts)
			if err != nil {
				res := api.ResponseApi(http.StatusInternalServerError, nil, err)
				c.JSON(http.StatusInternalServerError, res)
				return
			}
		}
		fileData := FileResponseData{
			FileName: filename,
			FileType: req[i].DirectoryFile,
		}

		resData = append(resData, fileData)

	}

	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusOK, res)
	return
}

func RollbackDeleteFile(ctx context.Context, m *MinioClient, resData []FileResponseData) {

	for _, v := range resData {

		opts := minio.RemoveObjectOptions{GovernanceBypass: true}
		err := m.mc.RemoveObject(ctx, m.bucketName, v.FileType+"/"+v.FileName, opts)
		if err != nil {
		}

	}
}
