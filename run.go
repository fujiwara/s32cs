package s32cs

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Run() {
	var reg *regexp.Regexp
	if r := os.Getenv("KEY_REGEXP"); r != "" {
		reg = regexp.MustCompile(r)
	} else {
		reg = nil
	}
	client := NewClient(session.New(), os.Getenv("ENDPOINT"), reg)

	handler := func(msg json.RawMessage) (interface{}, error) {
		var event SQSEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			return nil, err
		}
		if event.QueueURL != "" {
			log.Printf("starting process sqs %s", event.QueueURL)
			return nil, client.ProcessSQS(event.QueueURL)
		}
		var s3event S3Event
		if err := json.Unmarshal(msg, &s3event); err != nil {
			return nil, err
		}
		log.Println("starting process s3 event:", s3event.String())
		if err := client.Process(s3event); err != nil {
			return nil, err
		}
		return true, nil
	}

	env := os.Getenv("AWS_EXECUTION_ENV")
	runtimeAPI := os.Getenv("AWS_LAMBDA_RUNTIME_API")

	if strings.HasPrefix(env, "AWS_Lambda") || runtimeAPI != "" {
		lambda.Start(handler)
	} else if strings.HasPrefix(env, "Test_AWS_Lambda_go") {
		wrappedHandler := lambda.NewHandler(handler)
		_, err := wrappedHandler.Invoke(context.Background(), []byte("{}"))
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("[error] Lambda execution environment %s is not supported", env)
	}
}
