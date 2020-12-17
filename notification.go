package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Notification struct {
	Timestamp int64  `json:"timestamp"`
	Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
	Key       string `json:"key"`
}

func NewNotification(endpoint string, region string, bucket string, key string) *Notification {
	return &Notification{
		Timestamp: time.Now().Unix(),
		Endpoint:  endpoint,
		Region:    region,
		Bucket:    bucket,
		Key:       key,
	}
}

func getQueueURL(sess *session.Session, queueName *string) (*sqs.GetQueueUrlOutput, error) {
	svc := sqs.New(sess)
	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queueName,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func sendMessage(sess *session.Session, queueName string, groupId string, delaySeconds int64, notification *Notification) error {
	svc := sqs.New(sess)

	queueUrl, err := getQueueURL(sess, &queueName)
	if err != nil {
		return err
	}

	notification.Timestamp = time.Now().Unix()
	msg, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:   aws.Int64(delaySeconds),
		MessageGroupId: aws.String(groupId),
		MessageBody:    aws.String(string(msg)),
		QueueUrl:       queueUrl.QueueUrl,
	})

	return err
}
