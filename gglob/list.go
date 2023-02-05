package gglob

import "strings"

func listExpand(s string) (list []string) {
	last := len(s) - 1
	if len(s) >= 3 && s[0] == '{' && s[last] == '}' {
		s = s[1:last]
		list = strings.Split(s, ",")
	}
	return
}
