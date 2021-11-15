package connection_errors

import "fmt"

type DNSError struct {
	Host string
}

func (e *DNSError) Error() string {
	return fmt.Sprintf("%s - DNS lookup failed", e.Host)
}
