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

var (
	stringsTaggedTermOp = []string{"none", "=", "=~", "!=", "!=~"}
)

func (t TaggedTermOp) String() string {
	return stringsTaggedTermOp[t]
}

type TaggedTerm struct {
	Key         string
	Op          TaggedTermOp
	Value       string
	HasWildcard bool           // only for TaggedTermEq
	Glob        *WildcardItems // glob macher if HasWildcard
	Re          *regexp.Regexp // regexp
}

// Build compile regexp/glob
func (term *TaggedTerm) Build() (err error) {
	if term.HasWildcard {
		term.Glob = new(WildcardItems)
		if err = term.Glob.Parse(term.Value); err != nil {
			return
		}
		if len(term.Glob.Inners) == 0 {
			term.Value = term.Glob.P
			term.Glob = nil
			term.HasWildcard = false
		}
	} else if term.Op == TaggedTermMatch || term.Op == TaggedTermNotMatch {
		term.Re, err = regexp.Compile(term.Value)
		if err != nil {
			err = ErrExprInvalid{term.Value}
		}
	}
	return
}

func (term *TaggedTerm) Match(v string) bool {
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

// Build compile all regexp/glob
func (t TaggedTermList) Build() (err error) {
	for i := range t {
		if err = t[i].Build(); err != nil {
			return
		}
	}
	return
}

func (t TaggedTermList) Write(sb *strings.Builder) string {
	sb.WriteString("seriesByTag(")
	for i, term := range t {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('\'')
		sb.WriteString(term.Key)
		sb.WriteString(term.Op.String())
		sb.WriteString(term.Value)
		sb.WriteByte('\'')
	}
	sb.WriteString(")")
	return sb.String()
}

// MatchByPath match against tags map
func (terms TaggedTermList) MatchByTagsMap(tags map[string]string) bool {
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

// MatchByPath match against tags slice
func (terms TaggedTermList) MatchByTags(tags []Tag) bool {
	var i int
LOOP:
	for _, term := range terms {
		if len(tags) == i {
			if term.Op == TaggedTermNe || term.Op == TaggedTermNotMatch {
				// != and ~=! can be skiped and key can not exist
				continue LOOP
			}
			return false
		}
		if term.Key != tags[i].Key {
			// scan for tag
			for ; i < len(tags); i++ {
				if tags[i].Key == term.Key {
					break
				}
			}
			if len(tags) == i {
				if term.Op == TaggedTermNe || term.Op == TaggedTermNotMatch {
					// != and ~=! can be skiped and key can not exist
					continue LOOP
				}
				return false
			}
		}
		if !term.Match(tags[i].Value) {
			return false
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
				// terms[i].Glob = new(WildcardItems)
				// if err := terms[i].Glob.Parse(terms[i].Value); err != nil {
				// 	return nil, err
				// }
				// terms[i].Build()
			}
		case TaggedTermMatch, TaggedTermNotMatch:
		// 	var err error
		// 	terms[i].Re, err = regexp.Compile(terms[i].Value)
		// 	if err != nil {
		// 		return nil, ErrExprInvalid{s}
		// 	}
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

// PathTags split GraphiteMergeTree path format (like name?a=v1&b=v2&c=v3) into Tag's map
func PathTagsMap(path string) (tags map[string]string, err error) {
	name, args, ok := strings.Cut(path, "?")
	if !ok || strings.Contains(name, "=") {
		err = ErrPathInvalid{"name", "not found"}
	}
	tags = make(map[string]string)
	tags["__name__"] = escape.Unescape(name)
	var (
		kv, k, v string
	)
	for args != "" {
		if k, args, ok = strings.Cut(args, "="); ok {
			v, args, _ = strings.Cut(args, "&")
			key := escape.Unescape(k)
			tags[key] = escape.Unescape(v)
		} else {
			err = ErrPathInvalid{kv, "not delimited with ="}
			break
		}
	}
	return
}

func NextTag(tags string) (tag, value, next string, found bool) {
	tag, next, found = strings.Cut(tags, "=")
	if found {
		value, next, _ = strings.Cut(next, "&")
	}
	return
}

type Tag struct {
	Key   string
	Value string
}

// PathTags split GraphiteMergeTree path format (like name?a=v1&b=v2&c=v3) into Tag's slice
func PathTags(path string) (tags []Tag, err error) {
	name, args, ok := strings.Cut(path, "?")
	if !ok || strings.Contains(name, "=") {
		err = ErrPathInvalid{"name", "not found"}
	}
	tagsCount := strings.Count(args, "&") + 2
	tags = make([]Tag, 0, tagsCount)
	var (
		kv, k, v string
	)
	tags = append(tags, Tag{Key: "__name__", Value: escape.Unescape(name)})
	for args != "" {
		if k, v, args, ok = NextTag(args); ok {
			tags = append(tags, Tag{Key: escape.Unescape(k), Value: escape.Unescape(v)})
		} else {
			err = ErrPathInvalid{kv, "not delimited with ="}
			break
		}
	}
	return
}
