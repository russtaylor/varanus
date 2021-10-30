package varanus

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

func CheckSiteAvailability(_ context.Context, message PubSubMessage) error {
	url := string(message.Data)
	port, err := stringToPort(message.Attributes["port"])
	if err != nil {
		panic("Invalid port specified. Err: " + err.Error())
	}
	address := fmt.Sprintf("%v:%v", url, port)
	log.Infof("Beginning request to %v:%d", url, port)
	_, err = tls.Dial("tcp", address, nil)
	if err != nil {
		panic("Server doesn't support SSL Certificate. Err: " + err.Error())
	}
	return nil
}

func stringToPort(portString string) (int, error) {
	port, err := strconv.Atoi(portString)
	if err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, errors.New("invalid port number")
	}
	return port, nil
}
