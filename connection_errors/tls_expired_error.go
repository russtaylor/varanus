package connection_errors

import (
	"fmt"
)

type TLSExpiredError struct {
	Host string
}

func (e *TLSExpiredError) Error() string {
	return fmt.Sprintf("Certificate for %v is expired or invalid", e.Host)
}
