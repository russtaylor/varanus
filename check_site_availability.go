package varanus

import (
	"context"
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type PubSubMessage struct {
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	UrlString string `json:"url"`
	Port int `json:"port"`
	Secure bool `json:"ssl"`
	URL url.URL
}

func CheckSiteAvailability(_ context.Context, message PubSubMessage) error {
	parsedAttrs, err := validateAttributes(message.Attributes)
	if err != nil {
		log.Errorf("Error parsing request attributes!")
		panic("Error parsing request attributes. Error: " + err.Error())
	}
	targetUrl := parsedAttrs.UrlString
	port := parsedAttrs.Port
	//secure := message.Attributes.Secure
	address := fmt.Sprintf("%v:%v", targetUrl, port)
	log.Infof("Beginning request to %v:%d", targetUrl, port)
	_, err = tls.Dial("tcp", address, nil)
	if err != nil {
		log.Infof("Error type: %v", err)
		panic("Server doesn't support SSL Certificate. Err: " + err.Error())
	}
	return nil
}

func validateAttributes(attrs Attributes) (Attributes, error) {
	rawUrl, err := url.Parse(attrs.UrlString)
	if err != nil {
		return attrs, err
	}

	parsedAttrs := Attributes{
		UrlString: attrs.UrlString,
		Port: attrs.Port,
		Secure: attrs.Secure,
		URL: *rawUrl,
	}
	return parsedAttrs, nil
}
