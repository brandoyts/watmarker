package s3

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/core/port"
)

// this is an internal interface for mocking aws sdk
//
//go:generate mockgen -package s3 -source=s3.go -destination=s3_mock.go
type client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type Configuration struct {
	Bucket          string
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	BaseEndpoint    string
}

type Instance struct {
	Bucket string
	client client
}

func New(config Configuration) *Instance {
	cfg := aws.Config{
		Region:       config.Region,
		Credentials:  credentials.NewStaticCredentialsProvider(config.AccessKeyId, config.SecretAccessKey, ""),
		BaseEndpoint: &config.BaseEndpoint,
	}

	client := s3.NewFromConfig(cfg)

	return &Instance{
		Bucket: config.Bucket,
		client: client,
	}
}

var _ port.ImageStorage = (*Instance)(nil)

func (i *Instance) Upload(ctx context.Context, filename string, fileData []byte) error {
	_, err := i.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &i.Bucket,
		Body:   bytes.NewReader(fileData),
		Key:    &filename,
	})

	return err
}
