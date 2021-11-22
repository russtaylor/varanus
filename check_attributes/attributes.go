package check_attributes

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net/mail"
	"net/url"
)

const HttpProtocol = "http"
const HttpsProtocol = "https"
const DefaultTimeout = 30

type Attributes struct {
	UrlString string `json:"url"`
	Port      int    `json:"port"`
	SSL       bool   `json:"ssl"`
	Timeout   int    `json:"timeout"`
	Email     string `json:"email"`
	URL       url.URL
}

func ValidateAttributes(attrs Attributes) (Attributes, error) {
	urlObject, err := url.Parse(attrs.UrlString)
	if err != nil {
		return attrs, err
	}

	err = validateEmailAddress(attrs.Email)
	if err != nil {
		log.Fatalf("Unable to validate email address: %s", attrs.Email)
		return attrs, nil
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

	parsedAttrs := Attributes{
		UrlString: attrs.UrlString,
		Port:      port,
		SSL:       ssl,
		Timeout:   timeout,
		URL:       *urlObject,
		Email:     attrs.Email,
	}
	return parsedAttrs, nil
}

func validateEmailAddress(emailAddress string) error {
	_, err := mail.ParseAddress(emailAddress)
	if err != nil {
		return err
	}
	return nil
}

func getPort(attrs Attributes) (int, error) {
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

func shouldUseSsl(attrs Attributes) bool {
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
