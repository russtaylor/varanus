package varanus

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"russt.io/varanus/check_attributes"
	"russt.io/varanus/connection_errors"
	"strings"
	"time"
)

func SendAlertEmail(attrs check_attributes.Attributes, err error) {
	mailgunKey := strings.TrimSpace(os.Getenv("MAILGUN_KEY"))
	emailDomain := strings.TrimSpace(os.Getenv("MAILGUN_DOMAIN"))
	emailSender := strings.TrimSpace(os.Getenv("VARANUS_SENDER_EMAIL"))

	template := "undefined"
	title := "Undefined Error"

	switch err.(type) {
	default:
		log.Errorf("Unknown error type: %v", reflect.TypeOf(err))
		return
	case *connection_errors.DNSError:
		template = "dns_failure"
		title = fmt.Sprintf("DNS Failure for %v", attrs.UrlString)
	case *connection_errors.SSLExpiredError:
		template = "ssl_error"
		title = fmt.Sprintf("SSL Error for %v", attrs.UrlString)
	}
	now := time.Now()
	alertTime := now.UTC().Format(time.UnixDate)

	mg := mailgun.NewMailgun(emailDomain, mailgunKey)

	sender := emailSender
	subject := title
	recipient := attrs.Email

	message := mg.NewMessage(sender, subject, "", recipient)
	message.SetTemplate(template)

	err = message.AddTemplateVariable("address", attrs.UrlString)
	if err != nil {
		log.Errorf("Error setting address variable: %v", err.Error())
	}

	err = message.AddTemplateVariable("hostname", attrs.URL.Hostname())
	if err != nil {
		log.Errorf("Error setting hostname variable: %v", err.Error())
	}

	err = message.AddTemplateVariable("time", alertTime)
	if err != nil {
		log.Errorf("Error setting address variable: %v", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)
	if err != nil {
		log.Fatalf("Error sending message: %v", err.Error())
	}
	log.Infof("Alert message sent. ID: %s, Response: %s", id, resp)
}
