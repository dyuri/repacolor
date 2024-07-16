package color

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/mazznoer/csscolorparser"
)

func ParseColor(cstr string, usefallback bool) (col RepaColor, err error) {
	col = NOCOLOR
	c, err := csscolorparser.Parse(cstr)

	if err == nil {
		col = MakeColor(c)
	} else if (usefallback) {
		hash := md5.Sum([]byte(cstr))
		strhash := hex.EncodeToString(hash[:])
		c, err = csscolorparser.Parse("#" + strhash[:6])
		if err == nil {
			col = MakeColor(c)
		}
	} else {
		err = errors.New("cannot parse color")
	}

	return
}

func GetName(col RepaColor) (string, bool) {
	return csscolorparser.Color{R: col.R, G: col.G, B: col.B, A: col.A}.Name()
}
