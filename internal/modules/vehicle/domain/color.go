package vehicledomain

import (
	"errors"
	"regexp"
)

type Color string

var (
	ErrInvalidColor = errors.New("color must be a valid hex code")
	colorRegex      = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
)

func (c Color) Validate() error {
	if !colorRegex.MatchString(string(c)) {
		return ErrInvalidColor
	}
	return nil
}
