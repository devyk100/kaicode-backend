package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"log"
)

var (
	MESSAGE_GROUP_USER   = "user"
	MESSAGE_GROUP_SERVER = "server"
)

func SendMessage(messageBody string, messageGroup string) error {
	_, err := SQSClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:               &SQSQueueURL,
		MessageBody:            aws.String(messageBody),
		MessageGroupId:         aws.String(messageGroup),
		MessageDeduplicationId: aws.String(uuid.New().String()),
	})
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
		return err
	}
	return nil
}
