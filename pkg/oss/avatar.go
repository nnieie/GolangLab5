package oss

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/utils"
)

type AvatarOSSCli struct {
	bucketName   string
	publicDomain string
	sf           *utils.Snowflake
}

func NewAvatarOSSCli(bucketName, publicDomain string, snowflake *utils.Snowflake) *AvatarOSSCli {
	return &AvatarOSSCli{
		bucketName:   bucketName,
		publicDomain: publicDomain,
		sf:           snowflake,
	}
}

func (c *AvatarOSSCli) UploadAvatar(objectKey string, reader io.Reader) (fileURL string, err error) {
	_, err = r2Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
		Body:   reader,
	})

	if err != nil {
		logger.Errorf("failed to upload file to R2: %v", err)
		return "", err
	}

	// 上传成功后，拼接出可公开访问的 URL
	fileURL = fmt.Sprintf("%s/%s", c.publicDomain, objectKey)

	return fileURL, nil
}

func (c *AvatarOSSCli) GenerateImgName() (string, error) {
	sfid, err := c.sf.Generate()
	if err != nil {
		logger.Errorf("snowflake generate err: %v", err)
		return "", err
	}

	return utils.I64ToStr(sfid), nil
}
