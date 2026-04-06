package devicedomain

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidDeviceName = errors.New("device name must match pattern TRQ-{number}")
	deviceNameRegex      = regexp.MustCompile(`^TRQ-\d+$`)
)

type DeviceName string

func (n DeviceName) Validate() error {
	if !deviceNameRegex.MatchString(string(n)) {
		return ErrInvalidDeviceName
	}
	return nil
}
