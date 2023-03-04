package gtags

import (
	"regexp"

	"github.com/google/go-cmp/cmp"
)

var (
	cmpTransform = cmp.Transformer("Re", func(in *regexp.Regexp) string {
		if in == nil {
			return "<nil>"
		}
		return in.String()
	})
)
