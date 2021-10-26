package varanus

import (
	"context"
	"crypto/tls"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func CheckSiteAvailability(ctx context.Context, message PubSubMessage) error {
	url := string(message.Data)
	_, err := tls.Dial("tcp", url, nil)
	if err != nil {
		panic("Server doesn't support SSL Certificate. Err: " + err.Error())
	}
	return nil
}
