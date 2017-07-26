package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	apex "github.com/apex/go-apex"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fujiwara/s32cs"
)

var client *s32cs.Client

func run(msg json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event s32cs.SQSEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return nil, err
	}
	if event.QueueURL != "" {
		log.Printf("starting process sqs %s", event.QueueURL)
		return nil, client.ProcessSQS(event.QueueURL)
	}
	var s3event s32cs.S3Event
	if err := json.Unmarshal(msg, &s3event); err != nil {
		return nil, err
	}
	log.Printf("starting process s3 %s", s3event)
	if err := client.Process(s3event); err != nil {
		return nil, err
	}
	return true, nil
}

func main() {
	var reg *regexp.Regexp
	if r := os.Getenv("KEY_REGEXP"); r != "" {
		reg = regexp.MustCompile(r)
	} else {
		reg = nil
	}
	client = s32cs.NewClient(session.New(), os.Getenv("ENDPOINT"), reg)
	apex.HandleFunc(run)
}
