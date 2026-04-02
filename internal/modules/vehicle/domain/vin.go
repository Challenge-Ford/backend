package vehicledomain

import (
	"errors"
	"regexp"
)

type VIN string

var (
	ErrInvalidVIN = errors.New("VIN must be 17 alphanumeric characters (I, O and Q are not allowed)")
	vinRegex      = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)
)

func (v VIN) Validate() error {
	if !vinRegex.MatchString(string(v)) {
		return ErrInvalidVIN
	}
	return nil
}
