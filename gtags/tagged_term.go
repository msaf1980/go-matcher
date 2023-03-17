package gtags

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/msaf1980/go-matcher/glob"
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
	Glob        *glob.Glob     // glob macher if HasWildcard
	Re          *regexp.Regexp // regexp
}

func (t TaggedTerm) WriteString(buf *strings.Builder) {
	buf.WriteString(t.Key)
	buf.WriteString(t.Op.String())
	buf.WriteString(t.Value)
}

func (t TaggedTerm) String() string {
	var buf strings.Builder
	buf.Grow(24)
	t.WriteString(&buf)
	return buf.String()
}

// build compile regexp/glob
func (term *TaggedTerm) build() (err error) {
	if term.Op == TaggedTermMatch || term.Op == TaggedTermNotMatch {
		term.Re, err = regexp.Compile(term.Value)
		if err != nil {
			err = ErrExprInvalid{term.Value}
		}
	} else if items.HasWildcard(term.Value) {
		term.HasWildcard = true
		if term.Glob, err = glob.Parse(term.Value); err != nil {
			return err
		}
		term.Value = term.Glob.Node
		if len(term.Glob.Items) == 0 {
			term.Glob = nil
			term.HasWildcard = false
		} else {
			term.HasWildcard = true
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
	default:
		// must be unreacheable
		panic(fmt.Errorf("invalid op : %d", term.Op))
	}
}

// TaggedTermList is parsed seriesByTag expression
type TaggedTermList []TaggedTerm

func (t TaggedTermList) WriteString(buf *strings.Builder) {
	buf.WriteString("seriesByTag(")
	for i := 0; i < len(t); i++ {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\'')
		t[i].WriteString(buf)
		buf.WriteByte('\'')
	}
	buf.WriteString(")")
}

func (t TaggedTermList) String() string {
	var buf strings.Builder
	buf.Grow(24 * len(t))
	t.WriteString(&buf)
	return buf.String()
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

func ParseSeriesByTag(query string) (terms TaggedTermList, err error) {
	var (
		n          int
		conditions [128]string
	)
	n, err = seriesByTagArgs(query, conditions[:])
	if err != nil {
		return
	}

	return ParseTaggedConditions(conditions[:n])
}

func ParseTaggedConditions(conditions []string) (terms TaggedTermList, err error) {
	if len(conditions) == 0 {
		return
	}
	terms = make(TaggedTermList, len(conditions))

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
		} else {
			err = ErrExprInvalid{s}
			return
		}

		if terms[i].Key == "name" {
			terms[i].Key = "__name__"
		}

		if err = terms[i].build(); err != nil {
			return
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
		} else if terms[j].Key == "__name__" {
			return false
		}
		return terms[i].Key < terms[j].Key
	})

	return terms, nil
}

// GraphitePathTagsMap split Graphite path format (like name;a=v1;b=v2;c=v3) into Tag's map
func GraphitePathTagsMap(path string) (tags map[string]string, err error) {
	name, args, ok := strings.Cut(path, ";")
	if !ok || strings.Contains(name, "=") {
		err = ErrPathInvalid{"name", "not found"}
	}
	tags = make(map[string]string)
	tags["__name__"] = escape.Unescape(name)
	var (
		kv, k, v string
	)
	for args != "" {
		if k, v, args, ok = GraphiteNextTag(args); ok {
			key := escape.Unescape(k)
			tags[key] = escape.Unescape(v)
		} else {
			err = ErrPathInvalid{kv, "not delimited with ="}
			break
		}
	}
	return
}

// GraphitePathTagsMapB split Graphite path format (like name;a=v1;b=v2;c=v3) into prealloc Tag's map
func GraphitePathTagsMapB(path string, tags map[string]string) (err error) {
	name, args, ok := strings.Cut(path, ";")
	if !ok || strings.Contains(name, "=") {
		err = ErrPathInvalid{"name", "not found"}
	}
	for key := range tags {
		delete(tags, key)
	}
	tags["__name__"] = escape.Unescape(name)
	var (
		kv, k, v string
	)
	for args != "" {
		if k, v, args, ok = GraphiteNextTag(args); ok {
			key := escape.Unescape(k)
			tags[key] = escape.Unescape(v)
		} else {
			err = ErrPathInvalid{kv, "not delimited with ="}
			break
		}
	}
	return
}

// GraphitePathTags split Graphite tagged path format (like name;a=v1;b=v2;c=v3) into Tag's slice
func GraphitePathTags(path string) (tags []Tag, err error) {
	name, args, ok := strings.Cut(path, ";")
	if !ok || strings.Contains(name, "=") {
		err = ErrPathInvalid{"name", "not found"}
	}
	tagsCount := strings.Count(args, ";") + 2
	tags = make([]Tag, 0, tagsCount)
	var (
		kv, k, v string
	)
	tags = append(tags, Tag{Key: "__name__", Value: escape.Unescape(name)})
	for args != "" {
		if k, v, args, ok = GraphiteNextTag(args); ok {
			tags = append(tags, Tag{Key: escape.Unescape(k), Value: escape.Unescape(v)})
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

func GraphiteNextTag(tags string) (tag, value, next string, found bool) {
	tag, next, found = strings.Cut(tags, "=")
	if found {
		value, next, _ = strings.Cut(next, ";")
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
		err = ErrPathInvalid{"name", "not found in " + name}
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

// PathTagsMap split GraphiteMergeTree path format (like name?a=v1&b=v2&c=v3) into Tag's map
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
		if k, v, args, ok = NextTag(args); ok {
			key := escape.Unescape(k)
			tags[key] = escape.Unescape(v)
		} else {
			err = ErrPathInvalid{kv, "not delimited with ="}
			break
		}
	}
	return
}

// PathTagsMapB split GraphiteMergeTree path format (like name?a=v1&b=v2&c=v3) into prealloc Tag's map
func PathTagsMapB(path string, tags map[string]string) (err error) {
	name, args, ok := strings.Cut(path, "?")
	if !ok || strings.Contains(name, "=") {
		err = ErrPathInvalid{"name", "not found"}
	}
	for key := range tags {
		delete(tags, key)
	}
	tags["__name__"] = escape.Unescape(name)
	var (
		kv, k, v string
	)
	for args != "" {
		if k, v, args, ok = NextTag(args); ok {
			key := escape.Unescape(k)
			tags[key] = escape.Unescape(v)
		} else {
			err = ErrPathInvalid{kv, "not delimited with ="}
			break
		}
	}
	return
}

// TagsMap convert tags slice to map
func TagsMap(tags []Tag) map[string]string {
	tagsMap := make(map[string]string)
	for i := 0; i < len(tags); i++ {
		tagsMap[tags[i].Key] = tags[i].Value
	}
	return tagsMap
}
