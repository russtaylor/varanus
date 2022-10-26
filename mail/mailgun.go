package mail

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func SendMailgun(emailRecipient string, emailTitle string, htmlContent string) error {
	mailgunKey := strings.TrimSpace(os.Getenv("MAILGUN_KEY"))
	emailDomain := strings.TrimSpace(os.Getenv("MAILGUN_DOMAIN"))
	emailSender := strings.TrimSpace(os.Getenv("VARANUS_SENDER_EMAIL"))

	mg := mailgun.NewMailgun(emailDomain, mailgunKey)

	sender := emailSender
	subject := emailTitle
	recipient := emailRecipient

	message := mg.NewMessage(sender, subject, "", recipient)
	message.SetHtml(htmlContent)
	log.Tracef("Set HTML email content to: %s", htmlContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)
	if err != nil {
		log.Fatalf("Error sending message: %v", err.Error())
		return err
	}
	log.Infof("Alert message sent. ID: %s, Response: %s", id, resp)
	return nil
}
