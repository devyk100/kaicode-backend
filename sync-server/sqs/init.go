package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"os"
)

var SQSClient *sqs.Client
var SQSQueueURL string

func InitSQSClient() {
	customResolver := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"",
	))
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-south-1"), // Change to your AWS region
		config.WithCredentialsProvider(customResolver),
	)
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}
	client := sqs.NewFromConfig(cfg)
	SQSClient = client
	SQSQueueURL = os.Getenv("SQS_QUEUE_URL")
	if SQSQueueURL == "" {
		fmt.Println("SQS_QUEUE_URL is not set!")
	}
}
