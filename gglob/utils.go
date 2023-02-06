package gglob

import (
	"strings"
	"unicode/utf8"
)

func splitString(s string, start int) (string, string) {
	return s[:start], s[start:]
}

func nextString(s string, last int) string {
	return s[last:]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func PathLevel(path string) (string, int) {
	if path == "" {
		return path, 0
	}

	if path[len(path)-1] == '.' {
		return path[:len(path)-1], strings.Count(path, ".")
	}

	return path, strings.Count(path, ".") + 1
}

func pathSplit(path string) []string {
	if path == "" {
		return []string{}
	}

	if path[len(path)-1] == '.' {
		return strings.Split(path[:len(path)-1], ".")
	}

	return strings.Split(path, ".")
}

func hasEmptyParts(parts []string) bool {
	for _, part := range parts {
		if part == "" {
			return true
		}
	}
	return false
}

func WildcardCount(target string) (n int) {
	for _, c := range target {
		switch c {
		case '[', '{', '*', '?':
			n++
		}
	}
	return
}

func HasWildcard(target string) bool {
	return strings.ContainsAny(target, "[]{}*?")
}

func IndexWildcard(target string) int {
	return strings.IndexAny(target, "[]{}*?")
}

func IndexLastWildcard(target string) int {
	return strings.LastIndexAny(target, "[]{}*?")
}

func IntersectGlobs(globs []string) string {
	if len(globs) == 0 {
		return ""
	}
	if len(globs) == 0 {
		if pos := IndexWildcard(globs[0]); pos == -1 {
			return globs[0]
		} else {
			return globs[0][:pos]
		}
	}
	pos := 0
	for {
		c0, n := utf8.DecodeRuneInString(globs[0][pos:])
		switch c0 {
		case utf8.RuneError, '[', ']', '{', '}', '*', '?':
			return globs[0][:pos]
		}
		for i := 1; i < len(globs); i++ {
			if pos > len(globs[i])-1 {
				// glob shortest item
				return globs[0][:pos]
			}
			c, _ := utf8.DecodeRuneInString(globs[0][pos:])
			if c0 != c {
				return globs[0][:pos]
			}
			switch c {
			case utf8.RuneError, '[', ']', '{', '}', '*', '?':
				return globs[0][:pos]
			}
		}
		pos += n
	}
}

func splitList(s string) ([]string, error) {
	m := make([]string, 0, 8)

	return m, nil
}
