package vehicledomain

import (
	"errors"
	"regexp"
)

type Plate string

var (
	ErrInvalidPlate    = errors.New("invalid plate format")
	plateOldRegex      = regexp.MustCompile(`^[A-Z]{3}\d{4}$`)
	plateMercosulRegex = regexp.MustCompile(`^[A-Z]{3}\d[A-Z]\d{2}$`)
)

func (p Plate) Validate() error {
	s := string(p)
	if !plateOldRegex.MatchString(s) && !plateMercosulRegex.MatchString(s) {
		return ErrInvalidPlate
	}
	return nil
}
