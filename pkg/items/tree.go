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

type Store interface {
	Store(index int)
}

type MinStore struct {
	N int
}

func NewMinStore() *MinStore {
	return &MinStore{-1}
}

func (s *MinStore) Init() {
	s.N = -1
}

func (s *MinStore) Store(index int) {
	if s.N < 0 || s.N > index {
		s.N = index
	}
}

type TreeItem struct {
	NodeItem

	Reverse bool // for suffix

	Terminated string // end of chain (resulting raw/normalized globs)
	TermIndex  int    // rule num of end of chain (resulting glob), can be used in specific cases

	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	Childs []*TreeItem // next possible parts slice
}

func LocateChildTreeItem(childs []*TreeItem, node string, reverse bool) *TreeItem {
	for _, child := range childs {
		if child.Node == node && child.Reverse == reverse {
			return child
		}
	}
	return nil
}

// Match check string against []NodeItems (parsed wildcards or simple regular expression)
//
// return
//
// @matched counter for matched globs
func (treeItem *TreeItem) Match(s string, globs *[]string, index *[]int, first Store) (matched int) {
	for _, child := range treeItem.Childs {
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
			if s == "" && len(child.Terminated) > 0 {
				if globs != nil {
					*globs = append(*globs, child.Terminated)
				}
				if index != nil {
					*index = append(*index, child.TermIndex)
				}
				if first != nil {
					first.Store(child.TermIndex)
				}
				matched++
			}
			for _, subChild := range child.Childs {
				if n, _ := subChild.match(s, globs, index, first); n > 0 {
					matched += n
				}
			}
		} else {
			if n, _ := child.match(s, globs, index, first); n > 0 {
				matched += n
			}
		}
	}
	return
}

// match check string against []NodeItems (parsed wildcards or simple regular expression)
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func (item *TreeItem) match(s string, globs *[]string, index *[]int, first Store) (matched int, abort bool) {
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
		if s == "" && len(item.Terminated) > 0 {
			if globs != nil {
				*globs = append(*globs, item.Terminated)
			}
			if index != nil {
				*index = append(*index, item.TermIndex)
			}
			if first != nil {
				first.Store(item.TermIndex)
			}
			matched++
		}
		for i := 0; i < len(item.Childs); i++ {
			if n, _ := item.Childs[i].match(s, globs, index, first); n > 0 {
				matched += n
			}
		}
	case FindList:
		list := item.Item.(ItemList)

		if s == "" && len(item.Terminated) > 0 {
			if globs != nil {
				*globs = append(*globs, item.Terminated)
			}
			if index != nil {
				*index = append(*index, item.TermIndex)
			}
			if first != nil {
				first.Store(item.TermIndex)
			}
			matched++
		}

		if list.IsOptional() {
			if s == "" && len(item.Terminated) > 0 {
				if globs != nil {
					*globs = append(*globs, item.Terminated)
				}
				if index != nil {
					*index = append(*index, item.TermIndex)
				}
				if first != nil {
					first.Store(item.TermIndex)
				}
				matched++
			}
			for i := 0; i < len(item.Childs); i++ {
				if n, _ := item.Childs[i].match(s, globs, index, first); n > 0 {
					matched += n
				}
			}
		}

		for i := 0; i < list.Len(); i++ {
			if offset = list.MatchN(s, i); offset == -1 {
				continue
			}
			s := s[offset:]

			if s == "" && len(item.Terminated) > 0 {
				if globs != nil {
					*globs = append(*globs, item.Terminated)
				}
				if index != nil {
					*index = append(*index, item.TermIndex)
				}
				if first != nil {
					first.Store(item.TermIndex)
				}
				matched++
			}
			for i := 0; i < len(item.Childs); i++ {
				if n, _ := item.Childs[i].match(s, globs, index, first); n > 0 {
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
			if len(item.Terminated) > 0 {
				if globs != nil {
					*globs = append(*globs, item.Terminated)
				}
				if index != nil {
					*index = append(*index, item.TermIndex)
				}
				if first != nil {
					first.Store(item.TermIndex)
				}
				matched++
			}
			return
		} else {
			for i := 0; i < len(item.Childs); i++ {
				if n, _ := item.Childs[i].matchStar(s, globs, index, first); n > 0 {
					matched += n
				}
			}
		}
	default:
		panic(fmt.Errorf("unsupported find flag: %d", flag))
	}
	return
}

// matchStar check string against TreeItem after Star item
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func (item *TreeItem) matchStar(s string, globs *[]string, index *[]int, first Store) (matched int, abortGready bool) {
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

			if s == "" && len(item.Childs) == 0 && len(item.Terminated) > 0 {
				if globs != nil {
					*globs = append(*globs, item.Terminated)
				}
				if index != nil {
					*index = append(*index, item.TermIndex)
				}
				if first != nil {
					first.Store(item.TermIndex)
				}
				matched++
			} else {
				for i := 0; i < len(item.Childs); i++ {
					if n, _ := item.Childs[i].match(s, globs, index, first); n > 0 {
						matched += n
					}
				}
			}
		case FindList:
			list := item.Item.(ItemList)

			if list.IsOptional() && optional {
				optional = false
				if len(item.Terminated) > 0 {
					if globs != nil {
						*globs = append(*globs, item.Terminated)
					}
					if index != nil {
						*index = append(*index, item.TermIndex)
					}
					if first != nil {
						first.Store(item.TermIndex)
					}
					matched++
				}
				// refactor with two path for exclude
				for i := 0; i < len(item.Childs); i++ {
					if n, _ := item.Childs[i].matchStar(s, globs, index, first); n > 0 {
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
						if len(item.Terminated) > 0 {
							if globs != nil {
								*globs = append(*globs, item.Terminated)
							}
							if index != nil {
								*index = append(*index, item.TermIndex)
							}
							if first != nil {
								first.Store(item.TermIndex)
							}
							matched++
						}
					}
					for i := 0; i < len(item.Childs); i++ {
						if n, _ := item.Childs[i].match(s, globs, index, first); n > 0 {
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

		case FindNotSupported:
			panic("not supported in match")
		case FindForwarded:
			panic("forwarded match")
		case FindStar:
			if len(item.Childs) == 0 && len(item.Terminated) > 0 {
				if globs != nil {
					*globs = append(*globs, item.Terminated)
				}
				if index != nil {
					*index = append(*index, item.TermIndex)
				}
				if first != nil {
					first.Store(item.TermIndex)
				}
				matched++
				return
			}
			for i := 0; i < len(item.Childs); i++ {
				if n, _ := item.Childs[i].matchStar(s, globs, index, first); n > 0 {
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
