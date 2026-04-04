package devicedomain

import (
	"errors"
	"regexp"
)

var deviceNameRegex = regexp.MustCompile(`^TRQ-\d+$`)

type DeviceName string

func (n DeviceName) Validate() error {
	if !deviceNameRegex.MatchString(string(n)) {
		return errors.New("device name must match pattern TRQ-{number}")
	}
	return nil
}
