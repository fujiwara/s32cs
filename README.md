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
    "KEY_REGEXP": "Regexp to extract an endpoint from key"
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

## Requirements

- Go
- [Apex](http://apex.run)

## LICENSE

The MIT License (MIT)

Copyright (c) 2017 FUJIWARA Shunichiro / (c) 2017 KAYAC Inc.
