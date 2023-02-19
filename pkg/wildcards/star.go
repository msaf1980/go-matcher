package wildcards

import (
	"strings"
	"unicode/utf8"
)

type ItemStar struct{}

func (item ItemStar) Strings() []string {
	return nil
}

func (item ItemStar) Type() (typ ItemType, s string, c rune) {
	return ItemTypeOther, "", utf8.RuneError
}

func (item ItemStar) Locate(part string, nextItems []InnerItem) (offset int, support bool, _ int) {
	return -1, false, 0
}

func (item ItemStar) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	if part == "" && len(nextItems) == 0 {
		return true
	}

	nextOffset := 1 // string skip optimization
LOOP:
	for ; part != ""; part = part[nextOffset:] {
		part := part           // avoid overwrite outer loop
		nextItems := nextItems // avoid overwrite outer loop
		nextOffset = 1
		if len(nextItems) > 0 {
			nextItem := nextItems[0]
			// typ, _, _, vals := nextItem.Type()
			// gready skip scan, speedup find
			if idx, support, skip := nextItem.Locate(part, nextItems[1:]); support {
				if idx == -1 {
					break LOOP
				}
				nextOffset = idx
				part = part[idx:]
				nextItems = nextItems[1+skip:]
				found = true
			} else if vals := nextItem.Strings(); len(vals) > 0 {
				for _, v := range vals {
					part := part // avoid overwrite outer loop
				LOOP_LIST:
					for part != "" {
						if pos := strings.Index(part, v); pos == -1 {
							break LOOP_LIST
						} else {
							part = part[pos+len(v):]
							if len(nextItems) > 0 {
								if found = nextItems[1].Match(part, nextParts, nextItems[2:]); found {
									break LOOP_LIST
								}
							} else if part != "" {
								break LOOP_LIST
							}
						}
					}
				}
				return
			}
		} else {
			// all of string matched to *
			found = true
			break LOOP
		}
		if len(nextItems) > 0 {
			if found = nextItems[0].Match(part, nextParts, nextItems[1:]); found {
				break LOOP
			}
		}
	}
	return
}
