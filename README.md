# Poxy Proxy

HTTP service with configured endpoints to which you can post content and it will be uploaded to an S3 bucket.

It can optionally post a notification to an SQS queue.

## Config

See `config-example.json` for what should be a self-explanatory example.

## API

POST /{endpoint-name}/{key}

If the endpoint exists, uploads the body to S3 with the configured prefix and the key from the URL.

Data is transferred exactly as is, no interpretation of headers or content is performed.

If configured it will then send a notification to the configured SQS queue.

The message placed on the queue is a JSON object like this:

```json
{
    "timestamp": 1608163318,
    "endpoint": "nodeinfo",
    "region": "eu-west-2",
    "bucket": "stut-nodeinfo",
    "key": "desktop/system.json"
}
```

## License

Public domain so do what you want, just don't come crying to me if stuff hits fans.
