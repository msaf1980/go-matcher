package items

import (
	"sort"
	"strings"
)

func ListExpand(s string) (list []string, failed bool) {
	last := len(s) - 1
	if len(s) > 1 && s[0] == '{' && s[last] == '}' {
		s = s[1:last]
		if s == "" {
			return
		}
		list = strings.Split(s, ",")
		if len(list) > 0 {
			sort.Strings(list)
			// remove empty string in-place
			// from start
			i := 0
			for ; i < len(list) && list[i] == ""; i++ {
			}
			if i != len(list) {
				list = list[i:]
			}
		}
	} else {
		failed = true
	}

	return
}
