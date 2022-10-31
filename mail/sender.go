package mail

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"russt.io/varanus/connection_errors"
	"russt.io/varanus/connectivity_check"
	"time"
)

func SendAlertEmail(attrs connectivity_check.Attributes, err error) {
	emailTitle := "Undefined Error"
	failureType := "Unknown Failure"
	failureMessage := "Generic Message"
	shortFailureType := "Unknown Failure"

	switch err.(type) {
	default:
		log.Errorf("Unknown error type: %v", reflect.TypeOf(err))
		return
	case *connection_errors.DNSError:
		emailTitle = fmt.Sprintf("DNS Failure for %v", attrs.URL.Hostname())
		failureType = "DNS Failure"
		failureMessage = "DNS Lookup Failed"
		shortFailureType = "DNS Lookup Failed"
	case *connection_errors.TLSExpiredError:
		emailTitle = fmt.Sprintf("SSL Error for %v", attrs.URL.Hostname())
		failureType = "SSL Certificate Expired"
		failureMessage = "Your SSL certificate has expired"
		shortFailureType = "SSL Expired"
	case *connection_errors.TLSExpiresWithinPeriodError:
		emailTitle = fmt.Sprintf("SSL Expires soon for %v", attrs.URL.Hostname())
		failureType = "SSL Certificate Expires Soon"
		failureMessage = "Your SSL certificate expires in less than 7 days"
		shortFailureType = "SSL Expires Soon"
	}
	now := time.Now()
	alertTime := now.UTC().Format(time.UnixDate)

	templateVars := TemplateVars{
		FormattedTime:    alertTime,
		Url:              attrs.UrlString,
		Hostname:         attrs.URL.Hostname(),
		FailureType:      failureType,
		FailureMessage:   failureMessage,
		ShortFailureType: shortFailureType,
	}

	message, err := CompileTemplate(templateVars)
	if err != nil {
		log.Fatalf("Unable to compile template! %v", err.Error())
	}

	err = SendMailgun(attrs.Email, emailTitle, message)
	if err != nil {
		log.Fatal("Unable to send email via Mailgun!")
	}
}
