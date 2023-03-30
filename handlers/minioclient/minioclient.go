package minioclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	endPoint   = "uat-s3.universityapp.net"
	useSSL     = true
	accessKey  = "FLBWWX6CMZ6LGFE85QGZ"
	secretKey  = "tGy7i12hLpRpaTUB3OIPDP2Fb52OHRM7Rrlvjdek"
	bucketName = "researcher"
)

type MinioClient struct {
	mc *minio.Client
}

type MinioResponse struct {
	Status       int    `json:"status"`
	FilePath     string `json:"filePath"`
	ErrorMessage string `json:"errorMessage"`
	FileUrl      string `json:"fileUrl"`
}

func MinioClientConnect() *MinioClient {
	minioClient, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &MinioClient{minioClient}
}

func (m *MinioClient) UploadFile(c *gin.Context) {
	// single file
	ctx := context.Background()
	file, _ := c.FormFile("file")
	directory := c.PostForm("directory")
	if directory == "" {
		directory = "upload"
	}
	timenow := time.Now()
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

	response := MinioResponse{
		Status:       http.StatusCreated,
		FilePath:     newName,
		ErrorMessage: "",
		FileUrl:      "",
	}
	if _, err := m.mc.PutObject(ctx, bucketName, newName, fileBuffer, file.Size, minio.PutObjectOptions{ContentType: mimeType}); err != nil {
		response.Status = http.StatusInternalServerError
		response.FilePath = ""
		response.ErrorMessage = err.Error()
		c.JSON(http.StatusInternalServerError, response)
	} else {
		c.JSON(http.StatusCreated, response)
	}

}

func (m *MinioClient) GetFile(c *gin.Context) {
	filename := c.Param("filename")
	directory := c.Param("directory")
	ctx := context.Background()
	if directory == "" {
		directory = "upload"
	}

	reqParams := make(url.Values)
	response := MinioResponse{
		Status:       http.StatusOK,
		FilePath:     directory + "/" + filename,
		ErrorMessage: "",
		FileUrl:      "",
	}
	if presignedURL, err := m.mc.PresignedGetObject(ctx, bucketName, directory+"/"+filename, time.Second*24*60*60, reqParams); err != nil {
		response.Status = http.StatusInternalServerError
		response.FilePath = ""
		response.ErrorMessage = err.Error()
		c.JSON(http.StatusInternalServerError, response)
	} else {
		res := fmt.Sprintf("%s", presignedURL)
		response.FileUrl = res
		c.JSON(http.StatusOK, response)
	}
}

func (m *MinioClient) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	directory := c.Param("directory")
	ctx := context.Background()
	if directory == "" {
		directory = "upload"
	}
	response := MinioResponse{
		Status:       http.StatusOK,
		FilePath:     "",
		ErrorMessage: "",
		FileUrl:      "",
	}

	opts := minio.RemoveObjectOptions{GovernanceBypass: true}
	err := m.mc.RemoveObject(ctx, bucketName, directory+"/"+filename, opts)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.ErrorMessage = err.Error()
		c.JSON(http.StatusInternalServerError, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}
