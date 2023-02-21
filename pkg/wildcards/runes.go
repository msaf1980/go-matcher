package wildcards

import (
	"sort"
	"strings"
	"unicode/utf8"
)

// RunesExpand expand runes like [a-z0]
func RunesExpand(s string) (ranges ItemRuneRanges, failed bool) {
	if len(s) > 1 && s[0] == '[' && s[len(s)-1] == ']' {
		ranges = make(ItemRuneRanges, 0, 4)
		s = s[1 : len(s)-1]
		if len(s) == 0 {
			return
		}
		start := utf8.RuneError
		isRange := false
		for i, c := range s {
			if s[i] == '-' {
				isRange = true
			} else if isRange {
				if start == utf8.RuneError {
					start = c
					isRange = false
				} else {
					ranges = append(ranges, RuneRange{First: start, Last: c})
					start = utf8.RuneError
					isRange = false
				}
			} else {
				if start != utf8.RuneError {
					ranges = append(ranges, RuneRange{First: start, Last: start})
				}
				start = c
			}
		}
		if start != utf8.RuneError {
			ranges = append(ranges, RuneRange{First: start, Last: start})
		}
		sort.Slice(ranges, func(i, j int) bool {
			if ranges[i].First == ranges[j].First {
				return ranges[i].Last < ranges[j].Last
			}
			return ranges[i].First < ranges[j].First
		})

		// merge ranges
		j := 0
		n := len(ranges)
		for i := 0; i < n; i++ {
			pos := i
			if pos < n-1 {
				for next := i + 1; next < n; next++ {
					if ranges[pos].First == ranges[next].First || ranges[pos].Last+1 >= ranges[next].First {
						if ranges[pos].Last < ranges[next].Last {
							// merge two ranges, like [1-3 1-4] [1-3 2-5] [1-3 4-5 5-7]
							ranges[pos].Last = ranges[next].Last
						}
						// skip to next, merged
						i++
						ranges[i].Last = 0
					} else {
						break
					}
				}
			}

			if pos != j {
				ranges[j] = ranges[pos]
				ranges[pos].Last = 0
			}
			if ranges[j].Last != 0 {
				j++
			}
		}

		ranges = ranges[:j]
	} else {
		failed = true
	}

	return
}

type RuneRange struct {
	First rune
	Last  rune
}

// ItemRuneRanges is a range of rune symbols: [a-crzA-Z] is a range of a-c r z A-Z
type ItemRuneRanges []RuneRange

func (item ItemRuneRanges) Type() (typ ItemType, s string, c rune) {
	return ItemTypeOther, "", utf8.RuneError
}

func (item ItemRuneRanges) WriteString(buf *strings.Builder) {
	buf.WriteRune('[')
	for i, r := range item {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune(r.First)
		if r.First != r.Last {
			buf.WriteRune('-')
			buf.WriteRune(r.Last)
		}
	}
	buf.WriteRune(']')
}

func (item ItemRuneRanges) Strings() []string {
	return nil
}

func (item ItemRuneRanges) Locate(part string, nextItems []InnerItem) (offset int, support bool, _ int) {
	support = true
	for i, c := range part {
		if item.matchRune(c) {
			offset = i + 1
			return
		}
	}
	offset = -1
	return
}

func (item ItemRuneRanges) matchRune(c rune) bool {
	if c >= item[0].First && c <= item[len(item)-1].Last {
		for i := range item {
			// TODO: may be binary search for many ranges set ?
			if len(item) == 1 {
				if c >= item[0].First && c <= item[0].Last {
					return true
				}
			} else if c >= item[i].First && c <= item[i].Last {
				return true
			}
		}
	}
	return false
}

func (item ItemRuneRanges) Match(part string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		if item.matchRune(c) {
			found = true
			part = part[n:]

			if len(nextItems) > 0 {
				found = nextItems[0].Match(part, nextItems[1:])
			} else if part != "" && len(nextItems) == 0 {
				found = false
			}
		}
	}
	if found {
		if len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}

type ItemRuneMap map[rune]struct{}

func (item ItemRuneMap) Type() (typ ItemType, s string, c rune) {
	return ItemTypeOther, "", utf8.RuneError
}

func (item ItemRuneMap) Strings() []string {
	return nil
}

func (item ItemRuneMap) Locate(part string, nextItems []InnerItem) (offset int, support bool, _ int) {
	support = true
	for i, c := range part {
		if _, ok := item[c]; ok {
			offset = i + 1
			return
		}
	}
	offset = -1
	return
}

func (item ItemRuneMap) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		if _, ok := item[c]; ok {
			found = true
			part = part[n:]
		}
	}
	if found {
		if len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}
