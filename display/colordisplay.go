package display

import (
	"fmt"
	"strings"

	"github.com/dyuri/repacolor/color"
)

func TextColorDetails(c color.RepaColor) string {
	nameStr, _ := color.GetName(c)

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n", nameStr, c.Hex(), c.RgbString(), c.HslString(), c.LabString(), c.LchString(), c.OkLabString(), c.OkLchString())
}

func MergeStringsVertically(a, b string, width int) string {
	la := strings.Split(a, "\n")
	lb := strings.Split(b, "\n")

  if width == 0 {
    maxwidth := 0

    for _, s := range la {
      if len(s) > maxwidth {
        maxwidth = len(s)
      }
    }
    width = maxwidth
  }

	if len(la) < len(lb) {
		la = append(la, make([]string, len(lb)-len(la))...)
	} else if len(lb) < len(la) {
		lb = append(lb, make([]string, len(la)-len(lb))...)
	}

	var sb strings.Builder
	for i := 0; i < len(la); i++ {
    padding := width - len(la[i])
    if padding < 0 {
      padding = 0
    }
		sb.WriteString(fmt.Sprintf("%s%s %s", la[i], strings.Repeat(" ", padding), lb[i]))
		if (i < len(la) - 1) {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

