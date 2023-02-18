package wildcards

import "unicode/utf8"

type ItemOne struct{}

func (item ItemOne) Type() (typ ItemType, s string, c rune) {
	return ItemTypeOther, "", utf8.RuneError
}

func (item ItemOne) Strings() []string {
	return nil
}

func (item ItemOne) Locate(part string) (offset int, support bool) {
	return -1, false
}

func (item ItemOne) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		found = true
		part = part[n:]
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
