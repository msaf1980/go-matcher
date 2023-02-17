package items

import (
	"sort"
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

type ItemRuneRanges []RuneRange

func (item ItemRuneRanges) IsRune() (rune, bool) {
	return utf8.RuneError, false
}

func (item ItemRuneRanges) IsString() (string, bool) {
	return "", false
}

func (ItemRuneRanges) CanString() bool {
	return false
}

func (item ItemRuneRanges) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		if c >= item[0].First && c <= item[len(item)-1].Last {
			for i := range item {
				// TODO: may be binary search for many ranges set ?
				if len(item) == 1 {
					if c >= item[0].First && c <= item[0].Last {
						found = true
						part = part[n:]
						break
					}
				} else if c >= item[i].First && c <= item[i].Last {
					found = true
					part = part[n:]
					break
				}
			}
		}
	}
	if found {
		if part != "" && len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextParts, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}

type ItemRuneMap map[rune]struct{}

func (item ItemRuneMap) IsRune() (rune, bool) {
	return utf8.RuneError, false
}

func (item ItemRuneMap) IsString() (string, bool) {
	return "", false
}

func (ItemRuneMap) CanString() bool {
	return false
}

func (item ItemRuneMap) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		if _, ok := item[c]; ok {
			found = true
			part = part[n:]
		}
	}
	if found {
		if part != "" && len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextParts, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 || part == "" && len(nextItems) > 0 {
			found = false
		}
	}
	return
}
