package tests

import (
	"regexp"
	"strings"
)

func GlobRegexp(g string) string {
	s := g
	s = strings.ReplaceAll(s, ".", `\.`)
	s = strings.ReplaceAll(s, "$", `\$`)
	s = strings.ReplaceAll(s, "{", "(")
	s = strings.ReplaceAll(s, "}", ")")
	s = strings.ReplaceAll(s, "?", `\?`)
	s = strings.ReplaceAll(s, ",", "|")
	s = strings.ReplaceAll(s, "*", ".*")
	return "^" + s + "$"
}

func BuildGlobRegexp(g string) *regexp.Regexp {
	return regexp.MustCompile(GlobRegexp(g))
}
