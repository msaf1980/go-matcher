package items

import "strings"

// TODO: merge Star and One type

type ItemStar struct{}

// func (ItemStar) Type() NodeType {
// 	return NodeStar
// }

func (item ItemStar) IsString() (string, bool) {
	return "", false
}

func (item ItemStar) Match(part string, nextParts string, nextItems []InnerItem, _ bool) (found bool, _ int) {
	if part == "" && len(nextItems) == 0 {
		found = true
		return
	}

	var nextOffset int // string skip optimization
LOOP:
	for ; part != ""; part = part[nextOffset:] {
		part := part // avoid overwrite outer loop
		nextOffset = 1
		// nextItems := nextItems // avoid overwrite outer loop
		if len(nextItems) > 0 {
			if s, ok := nextItems[0].IsString(); ok {
				// speedup NodeString find
				if idx := strings.Index(part, s); idx == -1 {
					// string not found, no need star scan
					break LOOP
				} else {
					idx += len(s)
					nextOffset = idx
					part = part[idx:]
					nextItems = nextItems[1:]
				}
				// TODO: may be other optimization: may be for list
			}
		} else {
			// all of string matched to *
			found = true
			break LOOP
		}
		// if found {
		if part != "" && len(nextItems) > 0 {
			var idx int
			found, idx = nextItems[0].Match(part, nextParts, nextItems[1:], false)
			if idx == -1 {
				found = false
				break LOOP
			} else if idx > 0 {
				// first next item found, but others not, shit nextItems for avoid scan
				nextItems = nextItems[1:]
				nextOffset = idx
			}
			// } else if len(nextItems) == 0 {
			// 	found = true
			// 	break LOOP
		}
	}
	return
}
