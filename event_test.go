package s32cs_test

import (
	"encoding/json"
	"testing"

	"github.com/fujiwara/s32cs"
)

var s3eventSrc = []byte(`{
   "Records":[
      {
         "eventVersion":"2.0",
         "eventSource":"aws:s3",
         "awsRegion":"us-east-1",
         "eventTime":"1970-01-01T00:00:00.000Z",
         "eventName":"ObjectCreated:Put",
         "userIdentity":{
            "principalId":"AIDAJDPLRKLG7UEXAMPLE"
         },
         "requestParameters":{
            "sourceIPAddress":"127.0.0.1"
         },
         "responseElements":{
            "x-amz-request-id":"C3D13FE58DE4C810",
            "x-amz-id-2":"FMyUVURIY8/IgAtTv8xRjskZQpcIZ9KG4V5Wp6S7S/JRWeUWerMUE5JgHvANOjpD"
         },
         "s3":{
            "s3SchemaVersion":"1.0",
            "configurationId":"testConfigRule",
            "bucket":{
               "name":"mybucket",
               "ownerIdentity":{
                  "principalId":"A3NL1KOZZKExample"
               },
               "arn":"arn:aws:s3:::mybucket"
            },
            "object":{
               "key":"HappyFace.jpg",
               "size":1024,
               "eTag":"d41d8cd98f00b204e9800998ecf8427e",
               "versionId":"096fKKXTRTtl3on89fVO.nfljtsv6qko",
               "sequencer":"0055AED6DCD90281E5"
            }
         }
      }
   ]
}`)

func TestS3Event(t *testing.T) {
	var event s32cs.S3Event
	err := json.Unmarshal(s3eventSrc, &event)
	if err != nil {
		t.Errorf("S3Event unmarshal failed: %s", err)
	}
	if len(event.Records) != 1 {
		t.Errorf("insufficient S3Event.Records. expected 1 got %d", len(event.Records))
	}
	s3 := event.Records[0].S3
	name, key := s3.Bucket.Name, s3.Object.Key
	if name != "mybucket" {
		t.Errorf("wrong bucket name: %s", name)
	}
	if key != "HappyFace.jpg" {
		t.Errorf("wrong key: %s", key)
	}
	t.Log(event.String())
}
