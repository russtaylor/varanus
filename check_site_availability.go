package varanus

import (
	"context"
	log "github.com/sirupsen/logrus"
	"russt.io/varanus/connectivity_check"
	"russt.io/varanus/mail"
)

type PubSubMessage struct {
	CheckAttributes connectivity_check.Attributes `json:"attributes"`
}

func CheckSiteAvailability(_ context.Context, message PubSubMessage) error {
	parsedAttrs, err := connectivity_check.ValidateAttributes(message.CheckAttributes)
	if err != nil {
		log.Errorf("Error parsing request attributes!")
		panic("Error parsing request attributes. Error: " + err.Error())
	}

	log.Infof("Validating DNS record for %v", parsedAttrs.UrlString)
	err = connectivity_check.CheckDns(parsedAttrs)
	if err != nil {
		mail.SendAlertEmail(parsedAttrs, err)
		return err
	}

	log.Infof("Checking connectivity to %v", parsedAttrs.UrlString)
	err = connectivity_check.CheckConnectivity(parsedAttrs)
	if err != nil {
		panic("Error connecting to site... Error: " + err.Error())
	}
	if parsedAttrs.SSL == true {
		log.Infof("Checking SSL connection to %v", parsedAttrs.UrlString)
		err = connectivity_check.CheckTLSConnection(parsedAttrs)
		if err != nil {
			mail.SendAlertEmail(parsedAttrs, err)
			return err
		}
	}
	log.Infof("Successfully connected to '%v'", parsedAttrs.UrlString)

	return nil
}
