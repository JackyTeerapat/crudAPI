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

	timenow := time.Now()
	newName := timenow.Format("20060102150405") + "-" + file.Filename

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
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.String(http.StatusCreated, "Successfully uploaded : "+newName)
	}

}

func (m *MinioClient) GetFile(c *gin.Context) {
	filename := c.Param("filename")
	ctx := context.Background()

	reqParams := make(url.Values)
	if presignedURL, err := m.mc.PresignedGetObject(ctx, bucketName, filename, time.Second*24*60*60, reqParams); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		res := fmt.Sprintf("%s", presignedURL)
		c.String(http.StatusOK, res)
	}
}
