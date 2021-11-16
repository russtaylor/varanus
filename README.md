# Varanus

A tool to monitor the status of a site/address/API via cloud functions (for now GCP, but with goals to expand it).

Emails are sent only via Mailgun for now (with plans to add more providers later).

## Development

### Environment Variables

You'll need to configure your local environment first. Create a `.env` file. To send emails, you need a Mailgun API key
and a domain that's properly configured with Mailgun:

```txt
MAILGUN_KEY=key-00000000000000000000000000000000
MAILGUN_DOMAIN=example.com
VARANUS_SENDER_EMAIL='Site Alert <alerts@example.com>'
```

### Running Locally

Before you run locally, you'll need to source the `.env` file to load the variables:

```shell
set -o allexport; source .env; set +o allexport
```

To start serving the function locally, use `go run cmd/main.go`. This will serve the function on your local machine at
port 8080.

To actually run the code, you need to make a request. The following `curl` command will execute the command, but you
should change the URL/port as necessary.

```shell
curl localhost:8080 -X POST -H "Content-Type: application/json" -d '{
        "context": {
          "eventId":"123456789",
          "timestamp":"1970-01-01T01:01:01.001Z",
          "eventType":"google.pubsub.topic.publish",
          "resource":{
            "service":"pubsub.googleapis.com",
            "name":"projects/sample-project/topics/varanus-test",
            "type":"type.googleapis.com/google.pubsub.v1.PubsubMessage"
          }
        },
        "data": {
          "@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
          "attributes": {
             "url":"https://expired.badssl.com",
             "email":"email@example.com"
          }
        }
      }'
```

## Deployment

TODO
