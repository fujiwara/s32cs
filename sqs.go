package s32cs

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/Songmu/retry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (d *Client) ProcessSQS(queueURL string) error {
	for {
		output, err := d.queue.ReceiveMessage(
			&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueURL),
				MaxNumberOfMessages: aws.Int64(10),
				WaitTimeSeconds:     aws.Int64(0),
			},
		)
		if err != nil {
			return err
		}
		if len(output.Messages) == 0 {
			log.Println("info\tno messages available")
			return nil
		}
		for _, msg := range output.Messages {
			dec := json.NewDecoder(strings.NewReader(*(msg.Body)))
			var event S3Event
			if err := dec.Decode(&event); err != nil {
				log.Println("warn\tdecode error", err, *(msg.MessageId), *(msg.Body))
				continue
			}
			log.Println("info\tprocessing message", *(msg.MessageId))
			if err := d.Process(event); err != nil {
				log.Println("warn\tprocessing failed", err)
				continue
			}
			err := retry.Retry(3, time.Second, func() error {
				log.Println("info\tdeleting message", *(msg.MessageId))
				_, err := d.queue.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueURL),
					ReceiptHandle: msg.ReceiptHandle,
				})
				return err
			})
			if err != nil {
				log.Println("error\tdelete message failed", err)
			}
			log.Println("info\tdone", *(msg.MessageId))
		}
	}
	return nil
}
