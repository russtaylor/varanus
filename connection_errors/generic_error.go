package connection_errors

import "fmt"

type GenericError struct {
	Host string
}

func (e *GenericError) Error() string {
	return fmt.Sprintf("%s - Unable to connect, unknown reason.", e.Host)
}
