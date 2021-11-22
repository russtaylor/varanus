package connectivity_check

import (
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"russt.io/varanus/connection_errors"
)

func CheckTLSConnection(attrs Attributes) error {
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
