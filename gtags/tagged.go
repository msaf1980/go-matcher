package gtags

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/msaf1980/go-matcher/pkg/escape"
	"github.com/msaf1980/go-matcher/pkg/items"
)

// Based on github.com/go-graphite/graphite-clickhouse/finder/tagged.go

type TaggedTermOp int

const (
	TaggedTermEq       TaggedTermOp = 1 // =
	TaggedTermMatch    TaggedTermOp = 2 // =~
	TaggedTermNe       TaggedTermOp = 3 // !=
	TaggedTermNotMatch TaggedTermOp = 4 // !=~
)

type TaggedTerm struct {
	Key         string
	Op          TaggedTermOp
	Value       string
	HasWildcard bool           // only for TaggedTermEq
	Glob        *WildcardItems // glob macher if HasWildcard
	Re          *regexp.Regexp // regexp
}

func (term *TaggedTerm) Merge() {
	if term.HasWildcard && len(term.Glob.Inners) == 0 {
		term.Value = term.Glob.P
		term.Glob = nil
		term.HasWildcard = false
	}
}

func (term TaggedTerm) Match(v string) bool {
	switch term.Op {
	case TaggedTermEq:
		if term.HasWildcard {
			return term.Glob.Match(v)
		} else {
			return v == term.Value
		}
	case TaggedTermNe:
		if term.HasWildcard {
			return !term.Glob.Match(v)
		} else {
			return !(v == term.Value)
		}
	case TaggedTermMatch:
		return term.Re.MatchString(v)
	case TaggedTermNotMatch:
		return !term.Re.MatchString(v)
	}
	// must be unreacheable
	panic(fmt.Errorf("invalid op : %d", term.Op))
}

type TaggedTermList []TaggedTerm

func (terms TaggedTermList) MatchByTags(tags map[string]string) bool {
	for _, term := range terms {
		if v, ok := tags[term.Key]; ok {
			if !term.Match(v) {
				return false
			}
		} else if term.Op != TaggedTermNe && term.Op != TaggedTermNotMatch {
			// keys with != or ~=! check can be not exist, but others not
			return false
		}
	}
	return true
}

func nextTag(tags string) (tag, value, next string, found bool) {
	tag, next, found = strings.Cut(tags, "=")
	if found {
		value, next, _ = strings.Cut(next, "&")
	}
	return
}

// MatchByPath match against GraphiteMergeTree path format (lime name?a=v1&b=v2&c=v3)
func (terms TaggedTermList) MatchByPath(path string) bool {
	var (
		tag, value string
	)
	name, tags, _ := strings.Cut(path, "?")
	tag, value, tags, _ = nextTag(tags)

LOOP:
	for _, term := range terms {
		if term.Key == "__name__" {
			if !term.Match(name) {
				return false
			}
		} else {
			if tag == "" {
				if term.Op == TaggedTermNe || term.Op == TaggedTermNotMatch {
					// != and ~=! can be skiped and key can not exist
					continue LOOP
				}
				return false
			}
			if term.Key != tag {
				// scan for tag
				for {
					tag, value, tags, _ = nextTag(tags)
					if tag == "" {
						if term.Op == TaggedTermNe || term.Op == TaggedTermNotMatch {
							// != and ~=! can be skiped and key can not exist
							continue LOOP
						}
						return false
					}
					if tag == term.Key {
						break
					}
				}
			}
			if !term.Match(value) {
				return false
			}
		}
	}
	return true
}

func parseString(s string) (string, string, error) {
	if s[0] != '\'' && s[0] != '"' {
		return "", "", ErrNodeNotTerminated{s}
	}

	match := s[0]

	s = s[1:]

	var i int
	for i < len(s) && s[i] != match {
		i++
	}

	if i == len(s) {
		return "", "", ErrNodeNotTerminated{s}
	}

	return s[:i], s[i+1:], nil
}

func seriesByTagArgs(query string, args []string) (n int, err error) {
	// trim spaces
	e := strings.Trim(query, " ")
	if !strings.HasPrefix(e, "seriesByTag(") {
		err = ErrQueryInvalid{query}
		return
	}
	if e[len(e)-1] != ')' {
		err = ErrQueryInvalid{query}
		return
	}
	e = e[12 : len(e)-1]

	for len(e) > 0 {
		var arg string
		if e[0] == '\'' || e[0] == '"' {
			if arg, e, err = parseString(e); err != nil {
				return
			}
			// skip empty arg
			if arg != "" {
				if n == len(args) {
					err = ErrExprOverflow{query}
					return
				}
				args[n] = arg
				n++
			}
		} else if e[0] == ' ' || e[0] == ',' {
			e = e[1:]
		} else {
			err = ErrNodeNotTerminated{e}
			return
		}
	}
	return
}

func ParseSeriesByTag(query string) (TaggedTermList, error) {
	var conditions [128]string
	n, err := seriesByTagArgs(query, conditions[:])
	if err != nil {
		return nil, err
	}

	if n < 1 {
		return nil, ErrQueryInvalid{query}
	}

	return ParseTaggedConditions(conditions[:n])
}

func ParseTaggedConditions(conditions []string) (TaggedTermList, error) {
	terms := make(TaggedTermList, len(conditions))

	for i := 0; i < len(conditions); i++ {
		s := conditions[i]

		pos := strings.IndexAny(s, "=~!")
		if pos < 1 {
			return nil, ErrExprInvalid{s}
		}
		terms[i].Key = strings.TrimSpace(s[:pos])
		s = s[pos:]
		if strings.HasPrefix(s, "=") {
			if strings.HasPrefix(s, "=~") {
				terms[i].Op = TaggedTermMatch
				terms[i].Value = strings.TrimSpace(s[2:])
			} else {
				terms[i].Op = TaggedTermEq
				terms[i].Value = strings.TrimSpace(s[1:])
			}
		} else if strings.HasPrefix(s, "!=") {
			if strings.HasPrefix(s, "!=~") {
				terms[i].Op = TaggedTermNotMatch
				terms[i].Value = strings.TrimSpace(s[3:])
			} else {
				terms[i].Op = TaggedTermNe
				terms[i].Value = strings.TrimSpace(s[2:])
			}
		}

		if terms[i].Key == "name" {
			terms[i].Key = "__name__"
		}
		switch terms[i].Op {
		case TaggedTermEq, TaggedTermNe:
			if items.HasWildcard(terms[i].Value) {
				terms[i].HasWildcard = true
				terms[i].Glob = new(WildcardItems)
				if err := terms[i].Glob.Parse(terms[i].Value); err != nil {
					return nil, err
				}
				terms[i].Merge()
			}
		case TaggedTermMatch, TaggedTermNotMatch:
			var err error
			terms[i].Re, err = regexp.Compile(terms[i].Value)
			if err != nil {
				return nil, ErrExprInvalid{s}
			}
		default:
			return nil, ErrExprInvalid{s}
		}
	}

	sort.Slice(terms, func(i, j int) bool {
		if terms[i].Key == terms[j].Key {
			if terms[i].Op == terms[j].Op {
				return terms[i].Value < terms[j].Value
			} else {
				return terms[i].Op < terms[j].Op
			}
		}

		if terms[i].Key == "__name__" {
			return true
		}
		return terms[i].Key < terms[j].Key
	})

	return terms, nil
}

func PathTagsMap(path string) (tags map[string]string, err error) {
	name, args, ok := strings.Cut(path, "?")
	if !ok || strings.Contains(name, "=") {
		err = ErrPathInvalid{"name", "not found"}
	}
	tags = make(map[string]string)
	tags["__name__"] = name
	var (
		kv string
	)
	for args != "" {
		kv, args, _ = strings.Cut(args, "&")
		if k, v, ok := strings.Cut(kv, "="); ok {
			key := escape.Unescape(k)
			tags[key] = escape.Unescape(v)
		} else {
			err = ErrPathInvalid{kv, "not delimited with ="}
			return
		}
	}
	return
}
