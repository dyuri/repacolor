package color

import (
	"errors"

	"github.com/mazznoer/csscolorparser"
)

func ParseColor(cstr string) (RepaColor, error) {
	c, err := csscolorparser.Parse(cstr)

	if err == nil {
		return MakeColor(c), nil
	}

	return NoColor, errors.New("cannot parse color")
}
