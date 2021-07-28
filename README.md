# s32cs

Amazon CloudSearch document uploader via S3 event notification.

## Usage

Example Lambda functions configuration.

```js
{
  "FunctionName": "s32cs",
  "Description": "Amazon Cloudsearch uploader via S3 event notification.",
  "Environment": {
    "Variables": {
      "ENDPOINT": "<YOUR CloudSearch document default endpoint",
      "KEY_REGEXP": "Regexp to extract an endpoint from S3 object key"
    }
  },
  "Handler": "s32cs",
  "MemorySize": 128,
  "Role": "<YOUR role ARN>",
  "Runtime": "provided.al2",
  "Timeout": 60
}
```

Configure your S3 bucket. Set an event notification to the Lambda function.


### Example

```js
{
  "FunctionName": "s32cs",
  "Description": ""
  "Environment": {
    "Variables": {
      "ENDPOINT": "example-nregueirgrehbuigre.ap-northeast-1.cloudsearch.amazonaws.com",
      "KEY_REGEXP": "example/(.+?)/"
    }
  },
  "Handler": "s32cs",
  "MemorySize": 128,
  "Role": "arn:aws:iam::xxxxxxxxxxxx:role/s32cs_lambda_function",
  "Runtime": "provided.al2",
  "Timeout": 60
}
```

This configuration works as below.

- S3
   1. Event notification invokes the function.
- Lambda
   1. Read notified objects from S3.
   1. Convert the object (line delimitered JSON) to SDF.
   1. Upload SDF to CloudSearch.
     - endpoint is determined by `ENDPOINT` environment value(default) or extract from a object key by `KEY_REGEXP`.

## Source object file format

Line delimitered JSON only.

```json
{"id":"123","type":"add","fields":{"foo":"bar","bar":["A","B"]}}
{"id":"123","type":"delete"}
```

id, type (add or delete) columns are required.

## Using with Dead Letter Queue(DLQ)

If Lambda use DLQ (SQS), s32cs also works by jobs from SQS.

When s32cs invoked with payload like `{"queue_url:""}`, s32cs will fetch jobs from the SQS and upload.

```
// eample payload
{"queue_url":"https://sqs.ap-northeast-1.amazonaws.com/xxxxxxx/upload"}
```

You can this invocation periodically by CloudWatch Events scheduled jobs.

## Requirements

- Go

## LICENSE

The MIT License (MIT)

Copyright (c) 2017 FUJIWARA Shunichiro / (c) 2017 KAYAC Inc.
