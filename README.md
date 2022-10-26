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

### Create your Secret

Before your app can send any emails, you'll need to create the secret for your Cloud Function to use (you can just
specify it as an environment variable, but that's not very secure). Create a new secret by [following Google's
docs](https://cloud.google.com/secret-manager/docs/create-secret). Keep track of the `name` you use, you'll need
it in a bit.

### Create your Pub/Sub Topic

You'll need to create a pub/sub topic that will trigger your function. You can re-use the same topic for multiple site
checks, so you can give it a generic name.

```shell
gcloud pubsub topics create <PUBSUB_TOPIC>
```

### Create your Function

Creating the function is fairly straightforward. Run the following (filling in values as necessary) from your local
clone of this repo:

```shell
gcloud functions deploy <FUNCTION_NAME> \
  --entry-point=CheckSiteAvailability \
  --runtime=go116 \
  --set-secrets 'MAILGUN_KEY=<SECRET_NAME>:<SECRET_VERSION>' \
  --set-env-vars PROJECT_ID=<PROJECT_ID> \
  --set-env-vars VARANUS_SENDER_EMAIL=<EMAIL_ADDRESS> \
  --trigger-topic=<PUBSUB_TOPIC> \
  --project=<PROJECT_NAME> \
  --source=.
```

Make sure that completes successfully.

Note: I'm looking into using Cloud Functions 2nd Gen, but haven't tested it yet. So if you go that route, YMMV.

### Test Triggering your Function

To perform an (almost) end-to-end test, we can manually create some messages in Pub/Sub to ensure that the function is
working as expected.

For example:

```shell
gcloud pubsub topics publish <PUBSUB_TOPIC> --attribute=email="<EMAIL>",url="https://expired.badssl.com"
```

Assuming everything is set up correctly, that should send an email to the address specified. If you don't receive an
email, check the logs for your cloud function in the Google Console, or run:

```shell
gcloud functions logs read <FUNCTION_NAME>
```

### Schedule Periodic Runs

Finally, to make sure this is running and actually notifies you when problems arise, you need to schedule its runs.

Adjust the cron schedule as needed. The following will run the check every 5 minutes:

```shell
gcloud scheduler jobs create pubsub <JOB_NAME> --schedule="*/5 * * * *" --topic=<PUBSUB_TOPIC> --attributes=email="<EMAIL>",url="<URL_TO_CHECK>"
```
