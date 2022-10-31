package connectivity_check

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"russt.io/varanus/connection_errors"
	"time"
)

func CheckConnectivity(attrs Attributes) error {
	timeout := time.Duration(attrs.Timeout) * time.Second
	hostnameWithPort := fmt.Sprintf("%v:%v", attrs.URL.Host, attrs.Port)
	_, err := net.DialTimeout("tcp", hostnameWithPort, timeout)
	if err != nil {
		log.Infof("Unable to establish connection with '%v'", hostnameWithPort)
		log.Infof("More detailed error: %v", err.Error())
		return &connection_errors.GenericError{
			Host: attrs.URL.Host,
		}
	}
	return nil
}
