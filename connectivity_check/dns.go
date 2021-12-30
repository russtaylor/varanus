package connectivity_check

import (
	log "github.com/sirupsen/logrus"
	"net"
	"russt.io/varanus/connection_errors"
)

func CheckDns(attrs Attributes) error {
	_, err := net.LookupIP(attrs.URL.Host)
	if err != nil {
		log.Infof("Unable to resolve DNS for %v", attrs.UrlString)
		return &connection_errors.DNSError{
			Host: attrs.URL.Host,
		}
	}
	return nil
}
