package connection_errors

import (
	"fmt"
)

type SSLExpiredError struct {
	Host string
}

func (e *SSLExpiredError) Error() string {
	return fmt.Sprintf("Certificate for %v is expired or invalid", e.Host)
}
