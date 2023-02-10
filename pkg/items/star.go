package items

type ItemStar struct{}

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
				var idx int
				if idx, found = v.Locate(part); !found {
					break LOOP
				} else {
					nextOffset += idx
					part = part[idx:]
					nextItems = nextItems[1:]
					found = true
				}
				// TODO: may be other optimization: may be for list
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
