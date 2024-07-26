package display

import (
	"fmt"
	"strings"

	"github.com/dyuri/repacolor/color"
)

func TextColorDetails(c color.RepaColor) string {
	nameStr, hasName := color.GetName(c)

	if hasName {
		nameStr = fmt.Sprintf("Name:  %s\n", nameStr)
	}

	return fmt.Sprintf("%sHex:   %s\nRGB:   %s\nHSL:   %s\nLAB:   %s\nLCH:   %s\nOKLAB: %s\nOKLCH: %s\n", nameStr, c.Hex(), c.RgbString(), c.HslString(), c.LabString(), c.LchString(), c.OkLabString(), c.OkLchString())
}

func MergeStringsVertically(a, b string) string {
	la := strings.Split(a, "\n")
	lb := strings.Split(b, "\n")

	if len(la) < len(lb) {
		la = append(la, make([]string, len(lb)-len(la))...)
	} else if len(lb) < len(la) {
		lb = append(lb, make([]string, len(la)-len(lb))...)
	}

	var sb strings.Builder
	for i := 0; i < len(la); i++ {
		sb.WriteString(fmt.Sprintf("%s %s", la[i], lb[i]))
		if (i < len(la) - 1) {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

