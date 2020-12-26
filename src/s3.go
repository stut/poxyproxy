package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type SQS struct {
	Region       string `json:"region"`
	Queue        string `json:"queue"`
	Group        string `json:"group"`
	DelaySeconds int64  `json:"delay_seconds"`
}

func (config Endpoint) processS3Request(endpointName string, key string, contentType string, reader *bufio.Reader) error {
	sessionConfig := &aws.Config{}
	if len(config.Endpoint) != 0 {
		// Assume S3-compatible server.
		if len(config.Region) == 0 {
			config.Region = "us-east-1"
		}
		sessionConfig = &aws.Config{
			Endpoint:         aws.String(config.Endpoint),
			Region:           aws.String(config.Region),
			DisableSSL:       aws.Bool(!strings.Contains(config.Endpoint, "https://")),
			S3ForcePathStyle: aws.Bool(true),
		}
	} else if len(config.Region) != 0 {
		// AWS with region from endpoint configuration.
		sessionConfig = &aws.Config{Region: aws.String(config.Region)}
	} else {
		// AWS with default region or region from environment configuration.
	}

	sess := session.Must(session.NewSession(sessionConfig))
	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String(fmt.Sprintf("%s%s", config.Prefix, key)),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	if len(config.SQS.Queue) > 0 {
		notification := NewNotification(endpointName, config.Region, config.Bucket, key)
		err = sendMessage(sess, config.SQS, notification)
		if err != nil {
			return fmt.Errorf("file uploaded, notification failed: %v", err)
		}
	}

	return nil
}
