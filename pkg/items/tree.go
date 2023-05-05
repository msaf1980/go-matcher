package items

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

var (
	ErrIndexDup  = errors.New("duplicate index")
	ErrGlobExist = errors.New("glob already exist")
)

// Match check string against []NodeItems (parsed wildcards or simple regular expression)
//
// return
//
// @matched counter for matched globs
func (item *TreeItem) Match(s string, store Store) (matched int) {
	for _, child := range item.Childs {
		if child.Reverse {
			if len(s) < child.Item.MinLen() {
				continue
			}
			offset, flag := child.Item.MatchLast(s)
			if flag != FindDone {
				continue
			}
			if offset == -1 {
				continue
			}
			s := s[:offset]
			if s == "" && child.Terminate {
				store.Store(child.Query, child.Index)
				matched++
			}
			for _, subChild := range child.Childs {
				if n, _ := subChild.match(s, store); n > 0 {
					matched += n
				}
			}
		} else {
			if n, _ := child.match(s, store); n > 0 {
				matched += n
			}
		}
	}
	return
}

// func (treeItem *TreeItem) append(globs *[]string, index *[]int, first Store) {
// 	if globs != nil {
// 		*globs = append(*globs, treeItem.Terminated)
// 	}
// 	if index != nil {
// 		*index = append(*index, treeItem.TermIndex)
// 	}
// 	if first != nil {
// 		first.Store(treeItem.TermIndex)
// 	}
// }

// match check string against *TreeItem (parsed wildcards or simple regular expression)
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func (item *TreeItem) match(s string, store Store) (matched int, abort bool) {
	if len(s) < item.Item.MinLen() {
		return
	}
	offset, flag := item.Item.Match(s)
	if offset == -1 {
		return
	}
	s = s[offset:]
	switch flag {
	case FindDone:
		if s == "" && item.Terminate {
			store.Store(item.Query, item.Index)
			matched++
		}
		for i := 0; i < len(item.Childs); i++ {
			if n, _ := item.Childs[i].match(s, store); n > 0 {
				matched += n
			}
		}
	case FindList:
		list := item.Item.(ItemList)

		// if s == "" && item.Terminate {
		// 	item.append(globs, index, first)
		// 	matched++
		// }

		if list.IsOptional() {
			if s == "" && item.Terminate {
				store.Store(item.Query, item.Index)
				matched++
			}
			for i := 0; i < len(item.Childs); i++ {
				if n, _ := item.Childs[i].match(s, store); n > 0 {
					matched += n
				}
			}
		}

		for i := 0; i < list.Len(); i++ {
			if offset = list.MatchN(s, i); offset == -1 {
				continue
			}
			s := s[offset:]

			if s == "" && item.Terminate {
				store.Store(item.Query, item.Index)
				matched++
			}
			for i := 0; i < len(item.Childs); i++ {
				if n, _ := item.Childs[i].match(s, store); n > 0 {
					matched += n
				}
			}
		}
	case FindGroup:
		group := item.Item.(*Group)

		if group.IsOptional() {

			if s == "" && item.Terminate {
				store.Store(item.Query, item.Index)
				matched++
			}
			for i := 0; i < len(item.Childs); i++ {
				if n, _ := item.Childs[i].match(s, store); n > 0 {
					matched += n
				}
			}
		}

		for i := 0; i < len(group.Vals); i++ {
			if v, ok := group.Vals[i].(*Chain); ok {
				if n, _ := item.matchItemsInTree(s, v.Items, store); n > 0 {
					matched += n
				}
			} else {
				if n, _ := item.matchItemsInTree(s, []Item{group.Vals[i]}, store); n > 0 {
					matched += n
				}
			}
		}
	case FindNotSupported:
		panic("not supported in match")
	case FindForwarded:
		panic("forwarded match")
	case FindStar:
		if len(item.Childs) == 0 {
			if item.Terminate {
				store.Store(item.Query, item.Index)
				matched++
			}
			return
		} else {
			for i := 0; i < len(item.Childs); i++ {
				if n, _ := item.Childs[i].matchStar(s, store); n > 0 {
					matched += n
				}
			}
		}
	default:
		panic(fmt.Errorf("unsupported find flag: %d", flag))
	}
	return
}

// matchStar check string against *TreeItem after Star item
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func (item *TreeItem) matchStar(s string, store Store) (matched int, abortGready bool) {
	var (
		offset, length int
		flag           FindFlag
	)

	optional := true // flag for avoid repeated scan of optional list (empty value)
	for {
		if len(s) < item.Item.MinLen() {
			abortGready = true
			return
		}

		if offset, length, flag = item.Item.Find(s); offset == -1 {
			abortGready = true
			return
		}
		s = s[offset:]

		switch flag {
		case FindDone:
			s := s[length:]

			if s == "" && len(item.Childs) == 0 && item.Terminate {
				store.Store(item.Query, item.Index)
				matched++
			} else {
				for i := 0; i < len(item.Childs); i++ {
					if n, _ := item.Childs[i].match(s, store); n > 0 {
						matched += n
					}
				}
			}
		case FindList:
			list := item.Item.(ItemList)

			if list.IsOptional() && optional {
				optional = false
				if item.Terminate {
					store.Store(item.Query, item.Index)
					matched++
				}
				// refactor with two path for exclude
				for i := 0; i < len(item.Childs); i++ {
					if n, _ := item.Childs[i].matchStar(s, store); n > 0 {
						matched += n
					}
				}
			}

			for {
				var offsetN int
				offset := -1
				for i := 0; i < list.Len(); i++ {
					if offsetN, length = list.FindN(s, i); offsetN == -1 {
						continue
					}
					if offset == -1 || offset < offsetN {
						offset = offsetN
					}
					s := s[offsetN+length:]

					if s == "" {
						if item.Terminate {
							store.Store(item.Query, item.Index)
							matched++
						}
					}
					for i := 0; i < len(item.Childs); i++ {
						if n, _ := item.Childs[i].match(s, store); n > 0 {
							matched += n
						}
					}
				}
				if offset == -1 {
					break
				}
				s = s[offset:]
				if _, length = utf8.DecodeRuneInString(s); length < 1 {
					break
				}
				s = s[length:]
			}
			if offset == -1 {
				break
			}
			s = s[offset:]
			if _, length = utf8.DecodeRuneInString(s); length < 1 {
				break
			}
			s = s[length:]

		case FindGroup:
			group := item.Item.(*Group)

			if group.IsOptional() {
				optional = false
				if item.Terminate {
					store.Store(item.Query, item.Index)
					matched++
				}
				// refactor with two path for exclude
				for i := 0; i < len(item.Childs); i++ {
					if n, _ := item.Childs[i].matchStar(s, store); n > 0 {
						matched += n
					}
				}
			}

			for i := 0; i < len(group.Vals); i++ {
				if v, ok := group.Vals[i].(*Chain); ok {
					if n, _ := item.matchStarItemsInTree(s, v.Items, store); n > 0 {
						matched += n
					}
				} else {
					if n, _ := item.matchStarItemsInTree(s, []Item{group.Vals[i]}, store); n > 0 {
						matched += n
					}
				}
			}
			return

		case FindNotSupported:
			panic("not supported in match")
		case FindForwarded:
			panic("forwarded match")
		case FindStar:
			if len(item.Childs) == 0 && item.Terminate {
				store.Store(item.Query, item.Index)
				matched++
				return
			}
			for i := 0; i < len(item.Childs); i++ {
				if n, _ := item.Childs[i].matchStar(s, store); n > 0 {
					matched += n
				}
			}

			return
		default:
			panic(fmt.Errorf("unsupported find flag: %d", flag))
		}

		if _, length = utf8.DecodeRuneInString(s); length < 1 {
			break
		}
		s = s[length:]
	}
	return
}

func (item *TreeItem) matchNextTreeItem(s string, store Store) (matched int) {
	if s == "" && item.Terminate {
		store.Store(item.Query, item.Index)
		matched++
	}
	for _, child := range item.Childs {
		if n, _ := child.match(s, store); n > 0 {
			matched += n
		}
	}
	return
}

// matchItemsInTree check string against []Item (parsed wildcards or simple regular expression)
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func (item *TreeItem) matchItemsInTree(s string, items []Item, store Store) (matched int, abortGready bool) {

	var (
		pos, offset int
		flag        FindFlag
	)
	for {
		if pos == len(items) {
			if n := item.matchNextTreeItem(s, store); n > 0 {
				matched += n
			}
			return
		}
		if len(s) < items[pos].MinLen() {
			return
		}
		if offset, flag = items[pos].Match(s); offset == -1 {
			return
		}
		s = s[offset:]
		switch flag {
		case FindDone:
			pos++
		case FindList:
			list := items[pos].(ItemList)
			pos := pos + 1

			if list.IsOptional() {
				s := s
				if len(items) == pos {
					if n := item.matchNextTreeItem(s, store); n > 0 {
						matched += n
					}
					return
				} else if n, _ := item.matchItemsInTree(s, items[pos:], store); n > 0 {
					matched += n
				}
			}

			for i := 0; i < list.Len(); i++ {
				if offset = list.MatchN(s, i); offset == -1 {
					continue
				}
				s := s[offset:]
				if len(items) == pos {
					if n := item.matchNextTreeItem(s, store); n > 0 {
						matched += n
					}
					return
				} else if n, _ := item.matchItemsInTree(s, items[pos:], store); n > 0 {
					matched += n
				}
			}

			return
		case FindGroup:
			group := items[pos].(*Group)
			pos := pos + 1

			if group.IsOptional() {
				if len(items) == pos {
					if n := item.matchNextTreeItem(s, store); n > 0 {
						matched += n
					}
					return
				} else if n, _ := item.matchItemsInTree(s, items[pos:], store); n > 0 {
					matched += n
				}
			}

			for i := 0; i < len(group.Vals); i++ {
				if v, ok := group.Vals[i].(*Chain); ok {
					if n, _ := item.matchItemsInTree(s, v.Items, store); n > 0 {
						matched += n
					}
				} else {
					if n, _ := item.matchItemsInTree(s, []Item{group.Vals[i]}, store); n > 0 {
						matched += n
					}
				}
			}
			return
		case FindNotSupported:
			panic("not supported in match")
		case FindForwarded:
			panic("forwarded match")
		case FindStar:
			pos := pos + 1
			if len(items) == pos {
				if n := item.matchStarNextTreeItem(s, store); n > 0 {
					matched += n
				}
				return
			} else if n, _ := item.matchStarItemsInTree(s, items[pos:], store); n > 0 {
				matched += n
			}
		default:
			panic(fmt.Errorf("unsupported find flag: %d", flag))
		}
	}
}

func (item *TreeItem) matchStarNextTreeItem(s string, store Store) (matched int) {
	if item.Terminate {
		store.Store(item.Query, item.Index)
		matched++
	}
	for _, child := range item.Childs {
		if n, _ := child.matchStar(s, store); n > 0 {
			matched += n
		}
	}
	return
}

// matchStarItemsInTree check string against []Item after Star item
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func (item *TreeItem) matchStarItemsInTree(s string, items []Item, store Store) (matched int, abortGready bool) {

	var (
		offset, length int
		flag           FindFlag
	)
	if len(items) == 0 {
		if n := item.matchStarNextTreeItem(s, store); n > 0 {
			matched += n
		}
		return
	}

	for {
		if len(s) < items[0].MinLen() {
			abortGready = true
			return
		}

		if offset, length, flag = items[0].Find(s); offset == -1 {
			abortGready = true
			return
		}
		s = s[offset:]

		switch flag {
		case FindDone:
			sub := s[length:]
			if n, _ := item.matchItemsInTree(sub, items[1:], store); n > 0 {
				matched += n
			}
		case FindList:
			list := items[0].(ItemList)

			if list.IsOptional() {
				s := s
				if len(items) == 1 {
					if n := item.matchStarNextTreeItem(s, store); n > 0 {
						matched += n
					}
				} else if n, _ := item.matchStarItemsInTree(s, items[1:], store); n > 0 {
					matched += n
				}
			}

			for {
				var offsetN int
				offset := -1
				for i := 0; i < list.Len(); i++ {
					if offsetN, length = list.FindN(s, i); offsetN == -1 {
						// TODO: skip on next scan ?
						continue
					}
					if offset == -1 || offset < offsetN {
						offset = offsetN
					}
					s := s[offsetN+length:]
					if len(items) == 1 {
						if n := item.matchNextTreeItem(s, store); n > 0 {
							matched += n
						}
					} else if n, _ := item.matchItemsInTree(s, items[1:], store); n > 0 {
						matched += n
					}
				}
				if offset == -1 {
					break
				}
				s = s[offset:]
				// shift to one rune from minimal offset
				if _, length = utf8.DecodeRuneInString(s); length < 1 {
					break
				}
				s = s[length:]
			}

			return
		case FindGroup:
			group := items[0].(*Group)

			if group.IsOptional() {
				s := s
				if len(items) == 1 {
					if n := item.matchNextTreeItem(s, store); n > 0 {
						matched += n
					}
				} else if n, _ := item.matchStarItemsInTree(s, items[1:], store); n > 0 {
					matched += n
				}
			}

			for i := 0; i < len(group.Vals); i++ {
				if v, ok := group.Vals[i].(*Chain); ok {
					if n, _ := item.matchStarItemsInTree(s, v.Items, store); n > 0 {
						matched += n
					}
				} else {
					if n, _ := item.matchStarItemsInTree(s, []Item{group.Vals[i]}, store); n > 0 {
						matched += n
					}
				}
			}
			return

		case FindNotSupported:
			panic("not supported in match")
		case FindForwarded:
			panic("forwarded match")
		case FindStar:
			if len(items) == 1 {
				if n := item.matchStarNextTreeItem(s, store); n > 0 {
					matched += n
				}
			} else if n, _ := item.matchStarItemsInTree(s, items[1:], store); n > 0 {
				matched += n
			}
			return

		default:
			panic(fmt.Errorf("unsupported find flag: %d", flag))
		}

		if _, length = utf8.DecodeRuneInString(s); length < 1 {
			break
		}
		s = s[length:]
	}
	return
}
