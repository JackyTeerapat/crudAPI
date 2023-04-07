package minioclient

import (
	"CRUD-API/api"
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	mc         *minio.Client
	bucketName string
}

type FileResponseData struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FilePath string `json:"file_path"`
}

type UploadedFile struct {
	FileName   string `json:"file_name"`
	FileBase64 string `json:"base64"`
}

func MinioClientConnect() *MinioClient {
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
	return &MinioClient{minioClient, bucketName}
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

	directory := c.Param("directory")
	if directory == "" {
		directory = "upload"
	}

	resData := []FileResponseData{}
	for _, file := range files {
		timenow := time.Now()
		id := timenow.Format("20060102150405") + "-" + file.Filename
		newName := directory + "/" + timenow.Format("20060102150405") + "-" + file.Filename

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

		if _, err := m.mc.PutObject(ctx, m.bucketName, newName, fileBuffer, file.Size, minio.PutObjectOptions{ContentType: mimeType}); err != nil {
			res := api.ResponseApi(http.StatusInternalServerError, nil, err)
			c.JSON(http.StatusInternalServerError, res)
			return
		} else {
			fileData := FileResponseData{
				FileId:   id,
				FileName: file.Filename,
				FileType: directory,
				FilePath: newName,
			}
			resData = append(resData, fileData)
		}
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
		id := timenow.Format("20060102150405") + "-" + v.FileName
		newName := directory + "/" + timenow.Format("20060102150405") + "-" + v.FileName

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
				FileId:   id,
				FileName: v.FileName,
				FileType: directory,
				FilePath: newName,
			}
			resData = append(resData, fileData)
		}
	}
	res := api.ResponseApi(http.StatusOK, resData, nil)
	c.JSON(http.StatusCreated, res)
	return
}

func (m *MinioClient) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	directory := c.Param("directory")
	ctx := context.Background()
	if directory == "" {
		directory = "upload"
	}

	opts := minio.RemoveObjectOptions{GovernanceBypass: true}
	err := m.mc.RemoveObject(ctx, m.bucketName, directory+"/"+filename, opts)
	if err != nil {
		res := api.ResponseApi(http.StatusBadRequest, nil, err)
		c.JSON(http.StatusBadRequest, res)
	} else {
		res := api.ResponseApi(http.StatusOK, "Success", nil)
		c.JSON(http.StatusOK, res)
	}
	return
}
