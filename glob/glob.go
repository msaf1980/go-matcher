package glob

import (
	"strings"

	"github.com/msaf1980/go-matcher/pkg/items"
)

// Glob is glob matcher
type Glob struct {
	Glob string // raw glob
	Node string // optimized glob or value string if len(Inners) == 0

	MinLen int // min bytes len
	MaxLen int // -1 for unlimited

	Prefix string // prefix
	Suffix string
	Vals   map[string]struct{} // one list

	Items []items.Item
}

func (g *Glob) WriteRandom(buf *strings.Builder) {
	if len(g.Items) == 0 {
		buf.WriteString(g.Node)
	} else {
		for i := 0; i < len(g.Items); i++ {
			g.Items[i].WriteRandom(buf)
		}
	}

}

func (g *Glob) String() string {
	return g.Node
}

func (g *Glob) Raw() string {
	return g.Glob
}

func (g *Glob) Match(s string) (matched bool) {
	if g.Node == "*" {
		matched = true
		return
	}
	if len(s) < g.MinLen {
		return
	}
	if g.MaxLen > 0 && len(s) > g.MaxLen {
		return
	}
	if len(g.Items) == 0 {
		matched = (g.Node == s)
	} else {
		if g.Prefix != "" {
			if !strings.HasPrefix(s, g.Prefix) {
				// prefix not match
				return
			}
			s = s[len(g.Prefix):]
		}
		if g.Suffix != "" {
			if !strings.HasSuffix(s, g.Suffix) {
				// suffix not match
				return
			}
			s = s[:len(s)-len(g.Suffix)]
		}
		if len(g.Vals) > 0 {
			// large list optimization
			_, matched = g.Vals[s]
		} else {
			matched = items.MatchItems(s, g.Items)
		}
	}

	return
}

func Parse(glob string) (g *Glob, err error) {
	g = &Glob{Glob: glob}
	pos := items.IndexWildcard(glob)
	if pos == -1 {
		g.Node = glob
		g.MinLen = len(glob)
		g.MaxLen = g.MinLen
	} else {
		if pos > 0 {
			g.Prefix = glob[:pos] // prefix
			glob = glob[pos:]
			g.MinLen = len(g.Prefix)
			g.MaxLen = len(g.Prefix)
		}
		end := items.IndexLastWildcard(glob)
		if end == 0 && glob[0] != '?' && glob[0] != '*' {
			err = items.ErrNodeUnclosed{glob}
			return
		}
		if end < len(glob)-1 {
			end++
			g.Suffix = glob[end:]
			glob = glob[:end]
			g.MinLen += len(g.Suffix)
			g.MaxLen += len(g.Suffix)
		}

		switch glob {
		case "*":
			g.Items = []items.Item{items.Star(0)}
			g.MaxLen = -1 // unlimited
		case "?":
			g.Items = []items.Item{items.Any(1)}
			g.MinLen++
			g.MaxLen = items.AddMaxLen(g.MaxLen, 4) // rune max len is 4 bytes
		default:
			var (
				inner items.Item
			)
			innerCount := items.WildcardCount(glob)
			inners := make([]items.Item, 0, innerCount)

			for glob != "" {
				inner, glob, err = items.NextWildcardItem(glob)
				if err != nil {
					return
				}
				if inner == nil {
					continue
				}

				g.MinLen += inner.MinLen()
				g.MaxLen = items.AddMaxLen(g.MaxLen, inner.MaxLen())
				// try to in-place merge
				last := len(inners) - 1
				switch v := inner.(type) {
				case *items.String:
					if last == -1 {
						if g.Prefix == "" {
							g.Prefix = v.S
						} else {
							g.Prefix += v.S
						}
					} else {
						switch vv := inners[last].(type) {
						case *items.String:
							vv.Add(v.S)
						case items.Rune:
							v.PrependRune(rune(vv))
							inners[last] = v
						case items.Byte:
							v.PrependByte(byte(vv))
							inners[last] = v
						default:
							inners = append(inners, inner)
						}
					}
				case items.Rune:
					c := rune(v)
					if last == -1 {
						if g.Prefix == "" {
							g.Prefix = string(c)
						} else {
							g.Prefix += string(c)
						}
					} else {
						switch vv := inners[last].(type) {
						case *items.String:
							vv.AddRune(c)
						case items.Rune:
							inners[last] = vv.AppendRune(c)
						case items.Byte:
							inners[last] = vv.AppendRune(c)
						default:
							inners = append(inners, inner)
						}
					}
				case items.Byte:
					c := byte(v)
					if last == -1 {
						if g.Prefix == "" {
							g.Prefix = string(c)
						} else {
							g.Prefix += string(c)
						}
					} else {
						switch vv := inners[last].(type) {
						case *items.String:
							vv.AddByte(c)
						case items.Rune:
							inners[last] = vv.AppendByte(c)
						case items.Byte:
							inners[last] = vv.AppendByte(c)
						default:
							inners = append(inners, inner)
						}
					}
				case items.Any:
					if last == -1 {
						inners = append(inners, inner)
					} else {
						switch vv := inners[last].(type) {
						case items.Any:
							vv += v
							inners[last] = vv
						case items.Star:
							vv += items.Star(v)
							inners[last] = vv
						default:
							inners = append(inners, inner)
						}
					}
				case items.Star:
					if last == -1 {
						inners = append(inners, inner)
					} else {
						switch vv := inners[last].(type) {
						case items.Any:
							v += items.Star(vv)
							inners[last] = v
						case items.Star:
							vv += v
							inners[last] = vv
						default:
							inners = append(inners, inner)
						}
					}
				default:
					inners = append(inners, inner)
				}
			}
			if len(inners) > 1 {
				last := len(inners) - 1
				switch vv := inners[last].(type) {
				case *items.String:
					if g.Suffix == "" {
						g.Suffix = string(vv.S)
					} else {
						g.Suffix = string(vv.S) + g.Suffix
					}
					inners = inners[:last]
				case items.Rune:
					g.Suffix = string(vv) + g.Suffix
					inners = inners[:last]
				case items.Byte:
					g.Suffix = string(vv) + g.Suffix
					inners = inners[:last]
				}
			}
			if len(inners) > 0 {
				g.Items = inners
			}
		}

		if len(g.Items) == 0 {
			if g.Suffix == "" {
				g.Node = g.Prefix
				g.Prefix = ""
			} else {
				if g.Prefix == "" {
					g.Node = g.Suffix
				} else {
					g.Node = g.Prefix + g.Suffix
					g.Prefix = ""
				}
				g.Suffix = ""
			}
		} else {
			// TODO: write optimized glob to Node
			var buf strings.Builder
			buf.Grow(len(g.Glob))
			buf.WriteString(g.Prefix)
			for i := 0; i < len(g.Items); i++ {
				g.Items[i].WriteString(&buf)
			}
			buf.WriteString(g.Suffix)
			g.Node = buf.String()
		}

		if len(g.Items) == 1 {
			if v, ok := g.Items[0].(items.ItemList); ok && v.Len() > 2 {
				g.Vals = v.Map()
			}
		}
	}

	return
}

func ParseMust(glob string) *Glob {
	if g, err := Parse(glob); err != nil {
		panic(err)
	} else {
		return g
	}
}
