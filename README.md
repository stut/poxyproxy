# Poxy Proxy

HTTP service with configured endpoints to which you can post content and it will be either uploaded to an S3 bucket and optionally notify an SQS queue, or post it to a Slack incoming webhook.

## Context

The reason this exists is to support simple, restricted proxying of specific requests from nodes behind a gateway, removing the need to spread sensitive configuration around an internal network and limiting the features available.

For S3 this means limiting the spread of credentials, and restricting the buckets and prefixes that can be targeted.

For Slack this means a specific incoming webhook along with the option to prevent the customisation of the message username, avatar, and destination.

## Config

See `config-example.json` for what should be comprehensive and self-explanatory examples.

See `config-required.json` for an example that only contains required fields.

## API

Simple API with only one URL. All other URLs will return a 404.

```
POST /{endpoint-name}/{key}
```

If the endpoint exists, it will be be processed according to the configuration.

## Endpoint types

### S3

The body will be uploaded to an S3 object with the configured prefix and the key from the URL.

Data is transferred exactly as is, no interpretation of headers or content is performed.

If configured it will then send a notification to the SQS queue.

The message placed on the queue is a JSON object like this:

```json
{
    "timestamp": 1608163318,
    "endpoint": "nodeinfo",
    "region": "eu-west-2",
    "bucket": "nodeinfo",
    "key": "desktop/system.json"
}
```

Note that the region currently comes from the configuration file so will be empty if it's configured through other means.

#### Targeting S3-compatible servers

As you'll see from the example configuration files this is supported but has not been extensively tested. Pull requests to complete this functionality are more than welcome.

### Slack

If the request body is plain text it will be sent to the Slack URL using the `defaults` in the configuration.

If the request body is a JSON object and the `allow_overrides` option is `true` the defaults in the configuration will be updated with the object contents before being sent. If the option is `false` a `403 Forbidden` will be returned.

## Extensibility

The code is arranged to easily support additional endpoint types.

## License

Public domain so do what you want, just don't come crying to me if stuff hits fans.
