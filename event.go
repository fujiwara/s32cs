package s32cs

import (
	"encoding/json"
	"net/url"
	"time"
)

type SQSEvent struct {
	QueueURL string `json:"queue_url"`
}

type S3Event struct {
	Records []S3EventRecord `json:"Records"`
}

type S3EventRecord struct {
	EventVersion string    `json:"eventVersion"`
	EventSource  string    `json:"eventSource"`
	AwsRegion    string    `json:"awsRegion"`
	EventTime    time.Time `json:"eventTime"`
	EventName    string    `json:"eventName"`
	UserIdentity struct {
		PrincipalID string `json:"principalId"`
	} `json:"userIdentity"`
	RequestParameters struct {
		SourceIPAddress string `json:"sourceIPAddress"`
	} `json:"requestParameters"`
	ResponseElements struct {
		XAmzRequestID string `json:"x-amz-request-id"`
		XAmzID2       string `json:"x-amz-id-2"`
	} `json:"responseElements"`
	S3 struct {
		S3SchemaVersion string `json:"s3SchemaVersion"`
		ConfigurationID string `json:"configurationId"`
		Bucket          struct {
			Name          string `json:"name"`
			OwnerIdentity struct {
				PrincipalID string `json:"principalId"`
			} `json:"ownerIdentity"`
			Arn string `json:"arn"`
		} `json:"bucket"`
		Object struct {
			Key       string `json:"key"`
			Size      int    `json:"size"`
			ETag      string `json:"eTag"`
			VersionID string `json:"versionId"`
			Sequencer string `json:"sequencer"`
		} `json:"object"`
	} `json:"s3"`
}

func (r S3EventRecord) Parse() (bucket, key string, err error) {
	bucket = r.S3.Bucket.Name
	key, err = url.PathUnescape(r.S3.Object.Key)
	return
}

func (e S3Event) String() string {
	b, _ := json.Marshal(&e)
	return string(b)
}
