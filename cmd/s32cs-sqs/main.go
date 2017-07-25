package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Songmu/retry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/fujiwara/s32cs"
)

var (
	sess     = session.New()
	queueURL string
	endpoint string
)

func main() {
	flag.StringVar(&queueURL, "queue", "", "SQS queue URL")
	flag.StringVar(&endpoint, "endpoint", "", "CloudSearch endpoint URL")
	flag.Parse()

	if queueURL == "" || endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}

	csSess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(endpoint),
	}))
	domain := s32cs.NewDomain(csSess, sess)

	err := run(context.Background(), domain)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

func run(ctx context.Context, domain *s32cs.Domain) error {
	queue := sqs.New(sess)
	for {
		output, err := queue.ReceiveMessageWithContext(
			ctx,
			&sqs.ReceiveMessageInput{
				QueueUrl:              aws.String(queueURL),
				MaxNumberOfMessages:   aws.Int64(10),
				VisibilityTimeout:     aws.Int64(60),
				WaitTimeSeconds:       aws.Int64(0),
				MessageAttributeNames: aws.StringSlice([]string{"All"}),
				AttributeNames:        aws.StringSlice([]string{"All"}),
			},
		)
		if err != nil {
			return err
		}
		if len(output.Messages) == 0 {
			log.Println("no messages available")
			return nil
		}
		for _, msg := range output.Messages {
			dec := json.NewDecoder(strings.NewReader(*(msg.Body)))
			var event s32cs.S3Event
			if err := dec.Decode(&event); err != nil {
				log.Println("decode error", err, *(msg.MessageId), *(msg.Body))
				continue
			}
			log.Println("processing message", *(msg.MessageId))
			if err := domain.Process(event); err != nil {
				log.Println("processing failed", err)
				continue
			}
			err := retry.Retry(3, time.Second, func() error {
				log.Println("deleting message", *(msg.MessageId))
				_, err := queue.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueURL),
					ReceiptHandle: msg.ReceiptHandle,
				})
				return err
			})
			if err != nil {
				log.Println("delete message failed", err)
			}
			log.Println("done", *(msg.MessageId))
		}
	}
	return nil
}
