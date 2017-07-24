package main

import (
	"encoding/json"
	"os"

	apex "github.com/apex/go-apex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fujiwara/s32cs"
)

var domain *s32cs.Domain

func run(msg json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event S3Event
	err := json.Unmarshal(msg, &event)
	if err != nil {
		return nil, err
	}
	err := domain.Process(event)
	if err != nil {
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
