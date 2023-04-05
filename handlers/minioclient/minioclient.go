package minioclient

import (
	"context"
	"log"
	"net/http"
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
	Status       int                `json:"status"`
	Description  string             `json:"description"`
	ErrorMessage string             `json:"error_message"`
	FileData     []FileResponseData `json:"data"`
}

type FileResponseData struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FilePath string `json:"file_path"`
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
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	files := form.File["files"]

	directory := c.Param("directory")
	if directory == "" {
		directory = "upload"
	}

	response := MinioResponse{
		Status:       http.StatusCreated,
		Description:  "",
		ErrorMessage: "",
	}

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

		if _, err := m.mc.PutObject(ctx, bucketName, newName, fileBuffer, file.Size, minio.PutObjectOptions{ContentType: mimeType}); err != nil {
			response.Status = http.StatusInternalServerError
			response.ErrorMessage = err.Error()
			c.JSON(http.StatusBadRequest, response)
			return
		} else {
			response.Description = "Success"
			fileData := FileResponseData{
				FileId:   id,
				FileName: file.Filename,
				FileType: directory,
				FilePath: newName,
			}
			response.FileData = append(response.FileData, fileData)
		}
	}
	c.JSON(http.StatusCreated, response)
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
		ErrorMessage: "",
	}

	opts := minio.RemoveObjectOptions{GovernanceBypass: true}
	err := m.mc.RemoveObject(ctx, bucketName, directory+"/"+filename, opts)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, response)
	} else {
		response.Description = "Success"
		c.JSON(http.StatusOK, response)
	}
}
