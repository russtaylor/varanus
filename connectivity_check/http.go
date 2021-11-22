package connectivity_check

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func CheckConnectivity(attrs Attributes) error {
	timeout := time.Duration(attrs.Timeout) * time.Second
	hostnameWithPort := fmt.Sprintf("%v:%v", attrs.URL.Host, attrs.Port)
	_, err := net.DialTimeout("tcp", hostnameWithPort, timeout)
	if err != nil {
		log.Infof("Unable to establish connection with '%v'", hostnameWithPort)
		return errors.New("couldn't connect to " + hostnameWithPort)
	}
	return nil
}
