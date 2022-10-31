package connectivity_check

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/rickb777/date/period"
	log "github.com/sirupsen/logrus"
	"russt.io/varanus/connection_errors"
	"time"
)

func CheckTLSConnection(attrs Attributes, expirationPeriod period.Period) error {
	hostnameWithPort := fmt.Sprintf("%v:%v", attrs.URL.Host, attrs.Port)
	log.Infof("Beginning SSL check for '%v'", hostnameWithPort)
	conn, err := tls.Dial("tcp", hostnameWithPort, &tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Infof("Found invalid SSL config for %+v", err)
		return &connection_errors.TLSExpiredError{
			Host: attrs.URL.Host,
		}
	}

	defer func(conn *tls.Conn) {
		_ = conn.Close()
	}(conn)

	if err := conn.Handshake(); err != nil {
		log.Errorf("Handshake failed when connecting to %s", attrs.URL.Host)
		return err
	}

	peerCertificates := conn.ConnectionState().PeerCertificates
	certs := make([]*x509.Certificate, 0, len(peerCertificates))
	for _, cert := range peerCertificates {
		if cert.IsCA {
			continue
		}
		certs = append(certs, cert)
	}

	expirationDuration, _ := expirationPeriod.Duration()
	now := time.Now()
	expirationAlertThreshold := time.Now().Add(expirationDuration)

	for _, cert := range certs {
		if now.After(cert.NotAfter) {
			log.Warnf("Certificate for %s has expired", attrs.URL.Host)
			return &connection_errors.TLSExpiredError{
				Host: attrs.URL.Host,
			}
		}
		if expirationAlertThreshold.After(cert.NotAfter) {
			log.Warnf("Certificate for %s expires soon", attrs.URL.Host)
			return &connection_errors.TLSExpiresWithinPeriodError{
				Host:   attrs.URL.Host,
				Period: expirationPeriod,
			}
		}
	}

	return err
}
