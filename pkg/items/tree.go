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

type TreeItem struct {
	NodeItem

	Terminated []string // end of chain (resulting raw/normalized globs)
	TermIndex  []int    // rule num of end of chain (resulting glob), can be used in specific cases

	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	Childs []*TreeItem // next possible parts slice
}

func LocateChildTreeItem(childs []*TreeItem, node string) *TreeItem {
	for _, child := range childs {
		if child.Node == node {
			return child
		}
	}
	return nil
}

func setMin(v *int, a []int) {
	for i := 0; i < len(a); i++ {
		if *v < 0 || *v > a[i] {
			*v = a[i]
		}
	}
}

// Match check string against []NodeItems (parsed wildcards or simple regular expression)
//
// return
//
// @matched counter for matched globs
func (treeItem *TreeItem) Match(s string, globs *[]string, index *[]int, first *int) (matched int) {
	for _, child := range treeItem.Childs {
		if n, _ := child.match(s, globs, index, first); n > 0 {
			matched += n
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
func (item *TreeItem) match(s string, globs *[]string, index *[]int, first *int) (matched int, abort bool) {
	if s == "" && len(item.Childs) == 0 {
		if len(item.Terminated) > 0 {
			if globs != nil {
				*globs = append(*globs, item.Terminated...)
			}
			if index != nil {
				*index = append(*index, item.TermIndex...)
			}
			if first != nil {
				setMin(first, item.TermIndex)
			}
			matched++
		}
		return
	}
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
		for _, child := range item.Childs {
			if n, _ := child.match(s, globs, index, first); n > 0 {
				matched += n
			}
		}
	case FindList:
		list := item.Item.(ItemList)

		if s == "" && len(item.Terminated) > 0 {
			if globs != nil {
				*globs = append(*globs, item.Terminated...)
			}
			if index != nil {
				*index = append(*index, item.TermIndex...)
			}
			if first != nil {
				setMin(first, item.TermIndex)
			}
			matched++
		}

		if list.IsOptional() {
			for _, child := range item.Childs {
				if n, _ := child.match(s, globs, index, first); n > 0 {
					matched += n
				}
			}
		}

		for i := 0; i < list.Len(); i++ {
			if offset = list.MatchN(s, i); offset == -1 {
				continue
			}
			s := s[offset:]

			if s == "" {
				matchedChild := len(item.Terminated) > 0
				if matchedChild {
					if globs != nil {
						*globs = append(*globs, item.Terminated...)
					}
					if index != nil {
						*index = append(*index, item.TermIndex...)
					}
					if first != nil {
						setMin(first, item.TermIndex)
					}
					matched++
				}
				continue
			} else {
				for _, child := range item.Childs {
					if n, _ := child.match(s, globs, index, first); n > 0 {
						matched += n
					}
				}
			}
		}
	case FindNotSupported:
		panic("not supported in match")
	case FindForwarded:
		panic("forwarded match")
	case FindStar:
		if len(item.Childs) == 0 {
			matchedChild := len(item.Terminated) > 0
			if matchedChild {
				if globs != nil {
					*globs = append(*globs, item.Terminated...)
				}
				if index != nil {
					*index = append(*index, item.TermIndex...)
				}
				matched++
			}
			return
		} else {
			for _, child := range item.Childs {
				if n, _ := child.matchStar(s, globs, index, first); n > 0 {
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
func (item *TreeItem) matchStar(s string, globs *[]string, index *[]int, first *int) (matched int, abortGready bool) {
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

			if s == "" && len(item.Childs) == 0 {
				if len(item.Terminated) > 0 {
					if globs != nil {
						*globs = append(*globs, item.Terminated...)
					}
					if index != nil {
						*index = append(*index, item.TermIndex...)
					}
					if first != nil {
						setMin(first, item.TermIndex)
					}
					matched++
				}
			} else {
				for _, child := range item.Childs {
					if n, _ := child.match(s, globs, index, first); n > 0 {
						matched += n
					}
				}
			}
		case FindList:
			list := item.Item.(ItemList)

			if list.IsOptional() && optional {
				optional = false
				if s == "" || len(item.Terminated) > 0 {
					if len(item.Terminated) > 0 {
						if globs != nil {
							*globs = append(*globs, item.Terminated...)
						}
						if index != nil {
							*index = append(*index, item.TermIndex...)
						}
						if first != nil {
							setMin(first, item.TermIndex)
						}
						matched++
					}
				}
				// refactor with two path for exclude
				for _, child := range item.Childs {
					if n, _ := child.matchStar(s, globs, index, first); n > 0 {
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
						matchedChild := len(item.Terminated) > 0
						if matchedChild {
							if globs != nil {
								*globs = append(*globs, item.Terminated...)
							}
							if index != nil {
								*index = append(*index, item.TermIndex...)
							}
							if first != nil {
								setMin(first, item.TermIndex)
							}
							matched++
						}
					} else {
						for _, child := range item.Childs {
							if n, _ := child.match(s, globs, index, first); n > 0 {
								matched += n
							}
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
			if s == "" && len(item.Childs) == 0 {
				matchedChild := len(item.Terminated) > 0
				if matchedChild {
					if globs != nil {
						*globs = append(*globs, item.Terminated...)
					}
					if index != nil {
						*index = append(*index, item.TermIndex...)
					}
					if first != nil {
						setMin(first, item.TermIndex)
					}
					matched++
				}
				return
			} else {
				for _, child := range item.Childs {
					if n, _ := child.matchStar(s, globs, index, first); n > 0 {
						matched += n
					}
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
