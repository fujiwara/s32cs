package main

import (
	"flag"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fujiwara/s32cs"
)

func main() {
	var (
		sess      = session.New()
		queueURL  string
		endpoint  string
		keyRegexp string
	)

	flag.StringVar(&queueURL, "queue", "", "SQS queue URL")
	flag.StringVar(&endpoint, "default-endpoint", "", "CloudSearch default endpoint URL")
	flag.StringVar(&keyRegexp, "key-regex", "", "Regexp to extract an endpoint from key")
	flag.Parse()

	if queueURL == "" || endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	var reg *regexp.Regexp
	if keyRegexp != "" {
		reg = regexp.MustCompile(keyRegexp)
	} else {
		reg = nil
	}

	client := s32cs.NewClient(sess, endpoint, reg)

	err := client.ProcessSQS(queueURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
