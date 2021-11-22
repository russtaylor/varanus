package varanus

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"russt.io/varanus/check_attributes"
	"russt.io/varanus/connection_errors"
	"time"
)

type PubSubMessage struct {
	CheckAttributes check_attributes.Attributes `json:"attributes"`
}

func CheckSiteAvailability(_ context.Context, message PubSubMessage) error {

	parsedAttrs, err := check_attributes.ValidateAttributes(message.CheckAttributes)
	if err != nil {
		log.Errorf("Error parsing request attributes!")
		panic("Error parsing request attributes. Error: " + err.Error())
	}

	log.Infof("Validating DNS record for %v", parsedAttrs.UrlString)
	err = checkDns(parsedAttrs)
	if err != nil {
		SendAlertEmail(parsedAttrs, err)
		return err
	}

	log.Infof("Checking connectivity to %v", parsedAttrs.UrlString)
	err = checkConnectivity(parsedAttrs)
	if err != nil {
		panic("Error connecting to site... Error: " + err.Error())
	}
	if parsedAttrs.SSL == true {
		log.Infof("Checking SSL connection to %v", parsedAttrs.UrlString)
		err = checkSSLConnection(parsedAttrs)
		if err != nil {
			SendAlertEmail(parsedAttrs, err)
			return err
		}
	}
	log.Infof("Successfully connected to '%v'", parsedAttrs.UrlString)

	return nil
}

func checkDns(attrs check_attributes.Attributes) error {
	_, err := net.LookupIP(attrs.URL.Host)
	if err != nil {
		log.Infof("Unable to resolve DNS for %v", attrs.UrlString)
		return &connection_errors.DNSError{
			Host: attrs.URL.Host,
		}
	}
	return nil
}

func checkConnectivity(attrs check_attributes.Attributes) error {
	timeout := time.Duration(attrs.Timeout) * time.Second
	hostnameWithPort := fmt.Sprintf("%v:%v", attrs.URL.Host, attrs.Port)
	_, err := net.DialTimeout("tcp", hostnameWithPort, timeout)
	if err != nil {
		log.Infof("Unable to establish connection with '%v'", hostnameWithPort)
		return errors.New("couldn't connect to " + hostnameWithPort)
	}
	return nil
}

func checkSSLConnection(attrs check_attributes.Attributes) error {
	hostnameWithPort := fmt.Sprintf("%v:%v", attrs.URL.Host, attrs.Port)
	log.Infof("Beginning SSL check for '%v'", hostnameWithPort)
	_, err := tls.Dial("tcp", hostnameWithPort, nil)
	if err != nil {
		log.Infof("Found invalid SSL config for %+v", err)
		return &connection_errors.SSLExpiredError{
			Host: attrs.URL.Host,
		}
	}
	return err
}
