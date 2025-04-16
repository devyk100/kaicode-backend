package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func DeleteMessage(ctx context.Context, receiptHandle *string) error {
	_, err := SQSClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &SQSQueueURL,
		ReceiptHandle: receiptHandle,
	})
	return err
}
