package main

import (
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
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

	client := s32cs.NewClient(sess, endpoint)

	err := client.ProcessSQS(queueURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
