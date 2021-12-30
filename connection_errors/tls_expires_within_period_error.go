package connection_errors

import (
	"fmt"
	"github.com/rickb777/date/period"
)

type TLSExpiresWithinPeriodError struct {
	Host   string
	Period period.Period
}

func (e *TLSExpiresWithinPeriodError) Error() string {
	return fmt.Sprintf("Certificate for %v expires within the specified period: %v", e.Host, e.Period.String())
}
