package varanus

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/url"
	"russt.io/varanus/connection_errors"
	"time"
)

const HttpProtocol = "http"
const HttpsProtocol = "https"
const DefaultTimeout = 30

type PubSubMessage struct {
	CheckAttributes CheckAttributes `json:"attributes"`
}

type CheckAttributes struct {
	UrlString string `json:"url"`
	Port      int    `json:"port"`
	SSL       bool   `json:"ssl"`
	Timeout   int    `json:"timeout"`
	URL       url.URL
}

func CheckSiteAvailability(_ context.Context, message PubSubMessage) error {

	parsedAttrs, err := validateAttributes(message.CheckAttributes)
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
			panic("Error validating SSL connection... Error: " + err.Error())
		}
	}
	log.Infof("Successfully connected to '%v'", parsedAttrs.UrlString)

	return nil
}

func checkDns(attrs CheckAttributes) error {
	_, err := net.LookupIP(attrs.URL.Host)
	if err != nil {
		return &connection_errors.DNSError{
			Host: attrs.URL.Host,
		}
	}
	return nil
}

func checkConnectivity(attrs CheckAttributes) error {
	timeout := time.Duration(attrs.Timeout) * time.Second
	hostnameWithPort := fmt.Sprintf("%v:%v", attrs.URL.Host, attrs.Port)
	_, err := net.DialTimeout("tcp", hostnameWithPort, timeout)
	if err != nil {
		log.Infof("Unable to establish connection with '%v'", hostnameWithPort)
		return errors.New("couldn't connect to " + hostnameWithPort)
	}
	return nil
}

func checkSSLConnection(attrs CheckAttributes) error {
	hostnameWithPort := fmt.Sprintf("%v:%v", attrs.URL.Host, attrs.Port)
	log.Infof("Beginning SSL check for '%v'", hostnameWithPort)
	_, err := tls.Dial("tcp", hostnameWithPort, nil)
	if err != nil {
		log.Infof("Error type: %v", err.Error())
	}
	return err
}

func validateAttributes(attrs CheckAttributes) (CheckAttributes, error) {
	urlObject, err := url.Parse(attrs.UrlString)
	if err != nil {
		return attrs, err
	}

	attrs.URL = *urlObject

	timeout := attrs.Timeout
	if attrs.Timeout == 0 {
		timeout = DefaultTimeout
	}

	ssl := shouldUseSsl(attrs)
	port, err := getPort(attrs)
	if err != nil {
		return attrs, err
	}

	parsedAttrs := CheckAttributes{
		UrlString: attrs.UrlString,
		Port:      port,
		SSL:       ssl,
		Timeout:   timeout,
		URL:       *urlObject,
	}
	return parsedAttrs, nil
}

func getPort(attrs CheckAttributes) (int, error) {
	if attrs.Port != 0 {
		if attrs.Port < 0 || attrs.Port > 65535 {
			return 0, errors.New("invalid port specified")
		}
		return attrs.Port, nil
	}
	// No port specified - attempt to ascertain the port from the Scheme
	if attrs.URL.Scheme == HttpProtocol {
		return 80, nil
	} else if attrs.URL.Scheme == HttpsProtocol {
		return 443, nil
	}
	return 0, errors.New("unable to determine port")
}

func shouldUseSsl(attrs CheckAttributes) bool {
	if attrs.SSL == true {
		return true
	}
	if attrs.URL.Scheme == HttpsProtocol {
		return true
	}
	if attrs.Port == 443 {
		return true
	}
	return false
}
