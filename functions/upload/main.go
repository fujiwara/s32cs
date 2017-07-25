package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	apex "github.com/apex/go-apex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fujiwara/s32cs"
)

var domain *s32cs.Domain

func run(msg json.RawMessage, ctx *apex.Context) (interface{}, error) {
	if os.Getenv("FAIL") != "" {
		err := errors.New("env FAIL defined. failing now")
		log.Println(err)
		return nil, err
	}
	var event s32cs.S3Event
	if err := json.Unmarshal(msg, &event); err != nil {
		log.Println(err)
		return nil, err
	}
	if err := domain.Process(event); err != nil {
		log.Println(err)
		return nil, err
	}
	return true, nil
}

func main() {
	csSess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("REGION")),
		Endpoint: aws.String(os.Getenv("ENDPOINT")),
	}))
	s3Sess := session.New()
	domain = s32cs.NewDomain(csSess, s3Sess)
	apex.HandleFunc(run)
}
