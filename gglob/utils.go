package gglob

import (
	"strings"
)

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
