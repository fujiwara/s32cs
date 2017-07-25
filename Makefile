cmd/s32cs-sqs/s32cs-sqs: cmd/s32cs-sqs/main.go *.go
	cd cmd/s32cs-sqs && go build

clean:
	rm -f cmd/s32cs-sqs/s32cs-sqs
