package oss

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	cfg "github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/pkg/logger"
)

var r2Client *s3.Client

func InitR2Client() {
	r2Client = NewR2Client()
}

// NewR2Client 创建S3客户端，用于连接Cloudflare R2
func NewR2Client() *s3.Client {
	if cfg.CFR2Config == nil {
		logger.Fatalf("R2 config is nil: ensure config.Init(...) or LoadR2ConfigFromEnv() is called before using OSS")
		return nil
	}
	if cfg.CFR2Config.AccessKeyID == "" || cfg.CFR2Config.SecretAccessKey == "" || cfg.CFR2Config.Endpoint == "" {
		logger.Fatalf("R2 config incomplete: endpoint/access key/secret must be set")
		return nil
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.CFR2Config.AccessKeyID, cfg.CFR2Config.SecretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		logger.Fatalf("failed to load AWS config: %v", err)
		return nil
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.CFR2Config.Endpoint)
	})

	return s3Client
}
