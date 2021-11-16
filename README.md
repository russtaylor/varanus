# Varanus

A tool to monitor the status of a site/address/API via cloud functions (for now GCP).

## Development

To start serving the function locally, use `go run cmd/main.go`. This will serve the function on your local machine at port 8080.

To actually run the code, you need to make a request. The following `curl` command will execute the command, but you should change the URL/port as necessary.

```shell
curl localhost:8080 -X POST -H "Content-Type: application/json" -d '{
        "context": {
          "eventId":"1144231683168617",
          "timestamp":"2020-05-06T07:33:34.556Z",
          "eventType":"google.pubsub.topic.publish",
          "resource":{
            "service":"pubsub.googleapis.com",
            "name":"projects/sample-project/topics/gcf-test",
            "type":"type.googleapis.com/google.pubsub.v1.PubsubMessage"
          }
        },
        "data": {
          "@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
          "attributes": {
             "url":"https://expired.badssl.com"
          }
        }
      }'
```

## Deployment 

TODO
