package main

import (
	"bufio"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (config Config) processRequest(endpointName string, key string, reader *bufio.Reader) error {
	endpoint, gotEndpoint := config[endpointName]
	if !gotEndpoint {
		return fmt.Errorf("404")
	}

	sessionConfig := &aws.Config{}
	if len(endpoint.Region) != 0 {
		sessionConfig = &aws.Config{Region: aws.String(endpoint.Region)}
	}

	sess := session.Must(session.NewSession(sessionConfig))
	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(endpoint.Bucket),
		Key:    aws.String(fmt.Sprintf("%s%s", endpoint.Prefix, key)),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	if len(endpoint.Queue) > 0 {
		notification := NewNotification(endpointName, endpoint.Region, endpoint.Bucket, key)
		err = sendMessage(sess, endpoint.Queue, endpoint.Group, endpoint.DelaySeconds, notification)
		if err != nil {
			return fmt.Errorf("file uploaded, notification failed: %v", err)
		}
	}

	return nil
}
