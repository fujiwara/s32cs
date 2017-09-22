# s32cs

Amazon CloudSearch document uploader via S3 event notification.

## Usage

```js
// project.json
{
  "name": "s32cs",
  "description": "Amazon Cloudsearch uploader via S3 event notification.",
  "memory": 128,
  "timeout": 60,
  "language": "go",
  "role": "<YOUR role ARN>",
  "environment": {
    "ENDPOINT": "<YOUR CloudSearch document default endpoint",
    "KEY_REGEXP": "Regexp to extract an endpoint from S3 object key"
  }
}
```

Deploy a function to Lambda.

```
$ apex deploy
```

Configure your S3 bucket. Set an event notification to the Lambda function.


### Example

```js
{
  "name": "s32cs",
  "description": "Amazon Cloudsearch uploader via S3 event notification.",
  "memory": 128,
  "timeout": 60,
  "language": "go",
  "role": "arn:aws:iam::xxxxxxxxxxxx:role/s32cs_lambda_function",
  "environment": {
    "ENDPOINT": "example-nregueirgrehbuigre.ap-northeast-1.cloudsearch.amazonaws.com",
    "KEY_REGEXP": "example/(.+?)/"
  }
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

When s32cs invoked with `{"queue_url:""}` json object, s32cs will fetch jobs from the SQS and upload.

```
$ echo '{"queue_url":"https://sqs.ap-northeast-1.amazonaws.com/xxxxxxx/upload"} | apex invoke upload
```

You can this invocation periodically by CloudWatch Events scheduled jobs.

## Requirements

- Go
- [Apex](http://apex.run)

## LICENSE

The MIT License (MIT)

Copyright (c) 2017 FUJIWARA Shunichiro / (c) 2017 KAYAC Inc.
