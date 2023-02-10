package items

import "strings"

type ItemStar struct{}

// func (ItemStar) Type() NodeType {
// 	return NodeStar
// }

func (item ItemStar) IsString() (string, bool) {
	return "", false
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
			// speedup NodeString find
			case ItemString:
				s := string(v)
				if idx := strings.Index(part, s); idx == -1 {
					// string not found, no need star scan
					break LOOP
				} else {
					nextOffset += idx
					idx += len(s)
					part = part[idx:]
					nextItems = nextItems[1:]
					found = true
				}
				// TODO: may be other optimization: may be for list
			}
		} else {
			// all of string matched to *
			part = ""
			found = true
		}
		if found {
			if part != "" && len(nextItems) > 0 {
				found = nextItems[0].Match(part, nextParts, nextItems[1:])
			} else if part != "" || len(nextItems) > 0 {
				found = false
			}
			if found {
				break LOOP
			}
		}
	}
	return
}
