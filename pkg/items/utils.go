package items

import (
	"strings"
	"unicode/utf8"
)

func SplitString(s string, start int) (string, string) {
	return s[:start], s[start:]
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

func PathSplit(path string) (parts []string) {
	if path == "" {
		return []string{}
	}

	if path[len(path)-1] == '.' {
		parts = strings.Split(path[:len(path)-1], ".")
	} else {
		parts = strings.Split(path, ".")
	}
	return
}

func HasEmptyParts(parts []string) bool {
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
