package s32cs_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fujiwara/s32cs"
)

var csSess = session.Must(session.NewSession(&aws.Config{
	Region:   aws.String(os.Getenv("AWS_REGION")),
	Endpoint: aws.String(os.Getenv("CS_ENDPOINT")),
}))
var s3Sess = session.New()
var domain = s32cs.NewDomain(csSess, s3Sess)

func TestUpload(t *testing.T) {
	f, err := os.Open(os.Getenv("TEST_SDF"))
	if err != nil {
		t.Skip("test sdf file is not specified by TEST_SDF env")
		return
	}
	if err := domain.Upload(f); err != nil {
		t.Error(err)
	}
}

func TestProcess(t *testing.T) {
	var event s32cs.S3Event
	err := json.Unmarshal([]byte(os.Getenv("TEST_EVENT_JSON")), &event)
	if err != nil {
		t.Skip("TEST_EVENT_JSON invalid", err)
	}
	if err := domain.Process(event); err != nil {
		t.Error(err)
	}
}
