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

func BuildGlobSRegexp(globs []string) (re []*regexp.Regexp) {
	re = make([]*regexp.Regexp, len(globs))
	for i := 0; i < len(globs); i++ {
		re[i] = BuildGlobRegexp(globs[i])
	}
	return
}
