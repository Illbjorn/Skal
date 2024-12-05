package emit

import (
	"strings"

	"github.com/illbjorn/fstr"
)

func pairs(tmpl string, pairs ...string) string {
	return strings.Trim(fstr.Pairs(tmpl, pairs...), "\n")
}
