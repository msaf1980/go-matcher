package items

import "unicode/utf8"

type ItemStar struct{}

func (item ItemStar) IsRune() (rune, bool) {
	return utf8.RuneError, false
}

func (item ItemStar) IsString() (string, bool) {
	return "", false
}

func (ItemStar) CanString() bool {
	return false
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
			switch v := nextItem.(type) {
			// gready string skip scan, speedup NodeString find
			case ItemString:
				var idx int
				if idx, found = v.Locate(part); !found {
					break LOOP
				} else {
					nextOffset += idx
					part = part[idx:]
					nextItems = nextItems[1:]
					found = true
				}
			case ItemRune:
				var idx int
				if idx, found = v.Locate(part); !found {
					break LOOP
				} else {
					nextOffset += idx
					part = part[idx:]
					nextItems = nextItems[1:]
					found = true
				}
			case *ItemList:
				if v.ValsMin > 0 {
					// gready list skip scan, speedup find any first rune
					if idx := v.LocateFirst(part); idx == -1 {
						break LOOP
					} else {
						nextOffset += idx
						part = part[idx:]
					}
				}
			}
		} else {
			// all of string matched to *
			found = true
			break LOOP
		}
		if part != "" && len(nextItems) > 0 {
			if found = nextItems[0].Match(part, nextParts, nextItems[1:]); found {
				break LOOP
			}
		} else if part != "" || len(nextItems) > 0 {
			found = false
			break LOOP
		}

	}
	return
}
