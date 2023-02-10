package items

import "strings"

type ItemString string

// func (ItemString) Type() NodeType {
// 	return NodeString
// }

func (item ItemString) IsString() (string, bool) {
	return string(item), true
}

func (item ItemString) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	s := string(item)
	if part == s {
		// full match
		found = true
		part = ""
	} else if strings.HasPrefix(part, s) {
		// strip prefix
		found = true
		part = part[len(s):]
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

func (item ItemString) Locate(part string) (offset int, found bool) {
	s := string(item)
	if offset = strings.Index(part, s); offset != -1 {
		offset += len(s)
		found = true
	}
	return
}
