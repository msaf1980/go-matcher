package items

import "strings"

type ItemString string

// func (ItemString) Type() NodeType {
// 	return NodeString
// }

func (item ItemString) IsString() (string, bool) {
	return string(item), true
}

func (item ItemString) Match(part string, nextParts string, nextItems []InnerItem, gready bool) (found bool, offset int) {
	s := string(item)
	if gready {
		if offset = strings.Index(part, s); offset == -1 {
			return
		} else {
			found = true
			offset += len(s)
			part = part[offset:]
		}
	} else if strings.HasPrefix(part, s) {
		// strip prefix
		found = true
		part = part[len(s):]
	}
	if found {
		if part != "" && len(nextItems) > 0 {
			found, _ = nextItems[0].Match(part, nextParts, nextItems[1:], false)
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}
