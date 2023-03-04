package gglob

import (
	"strings"

	"github.com/msaf1980/go-matcher/glob"
	"github.com/msaf1980/go-matcher/pkg/items"
)

// GGlob is glob matcher (dot-separated, like a.b*.c), writted for graphite project
type GGlob struct {
	Glob string // raw glob
	Node string // optimized glob or value string if len(Parts) == 0

	MinLen int // min bytes len
	MaxLen int // -1 for unlimited

	Parts []*glob.Glob
}

func (g *GGlob) Match(path string) (matched bool) {
	if path == "" {
		return
	}
	path, partsCount := PathLevel(path)
	if len(g.Parts) != partsCount {
		return
	}

	if len(path) < g.MinLen {
		return
	}
	if g.MaxLen > 0 && len(path) > g.MaxLen {
		return
	}

	var part string
	for i := 0; i < len(g.Parts); i++ {
		part, path, _ = strings.Cut(path, ".")
		if part == "" {
			return
		}

		if !g.Parts[i].Match(part) {
			return
		}
	}

	matched = path == ""
	return
}

func (g *GGlob) MatchByParts(parts []string, length int) (matched bool) {
	if len(parts) == 0 {
		return
	}
	if len(g.Parts) != len(parts) {
		return
	}
	if length > 0 {
		if length < g.MinLen {
			return
		}
		if g.MaxLen > 0 && length > g.MaxLen {
			return
		}
	}

	for i := 0; i < len(g.Parts); i++ {
		if !g.Parts[i].Match(parts[i]) {
			return
		}
	}

	matched = true
	return
}

func Parse(s string) (gg *GGlob, err error) {
	var (
		level int
		part  string
		g     *glob.Glob
	)

	s, level = PathLevel(s)

	gg = &GGlob{Glob: s, Parts: make([]*glob.Glob, 0, level)}

	nextParts := s
	for nextParts != "" {
		part, nextParts, _ = strings.Cut(nextParts, ".")
		if part == "" {
			err = items.ErrNodeEmpty{s}
			return
		}

		if g, err = glob.Parse(part); err != nil {
			return
		}

		gg.Parts = append(gg.Parts, g)
		gg.MinLen += g.MinLen
		gg.MaxLen = items.AddMaxLen(gg.MaxLen, g.MaxLen)
	}

	var buf strings.Builder
	buf.Grow(len(s))
	for i, g := range gg.Parts {
		if i > 0 {
			buf.WriteByte('.')
		}
		buf.WriteString(g.Node)
	}
	gg.Node = buf.String()

	return
}

func ParseMust(s string) *GGlob {
	if gg, err := Parse(s); err != nil {
		panic(err)
	} else {
		return gg
	}
}
