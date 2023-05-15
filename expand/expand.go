package expand

import (
	"errors"
	"io"
	"strings"

	"github.com/msaf1980/go-matcher/pkg/utils"
)

var (
	ErrNotClosed = errors.New("expression not closed")

	asciiSet = utils.MakeASCIISetMust("[]{}*?")
)

// Expand takes the string contains the shell expansion expression and returns list of strings after they are expanded (from begin).
//
// Argument max for restict max expanded results, > 0 - restuct  max expamnded results, 0 - disables expand, -1 - unlimited, etc.
func Expand(in string, max, depth int) ([]string, error) {
	if max == 0 {
		return []string{in}, nil
	}
	exps := parseExpr(in)
	if len(exps) == 1 && exps[0].typ == expString {
		return []string{exps[0].body}, nil
	}

	count := 1
	for i := 0; i < len(exps); i++ {
		count *= exps[i].count()
		if exps[i].typ == expWildcard {
			break
		}
	}

	result := make([]string, 0, count)

	var err error
	buf := make([]byte, 0, len(in))
	if depth > 0 {
		depth++
	}
	if result, _, err = expand(exps, result, 0, max, depth, buf); err != nil {
		return nil, err
	}

	return result, nil
}

// getPair returns the top level expression.
func getPair(in string) (start, stop int) {
	start = -1
	stop = -1
	for i, c := range in {
		switch c {
		case '*', '?':
			// break, no expand after star
			if start == -1 {
				start = i
			}
			return
		case '{', '[':
			if start == -1 {
				start = i
			} else {
				return
			}
		case '}':
			if start == -1 || in[start] == '[' {
				if start == -1 {
					start = i
				}
				return
			}
			stop = i
			return
		case ']':
			if start == -1 || in[start] == '{' {
				if start == -1 {
					start = i
				}
				return
			}
			stop = i
			return
		}
	}

	return
}

func parseExpr(in string) []expression {
	// var starBreak int
	start, stop := getPair(in)
	if stop == -1 {
		if start == -1 {
			return []expression{{body: in}}
		} else if start == 0 {
			pos := asciiSet.LastIndex(in) + 1
			if pos < len(in) {
				return []expression{
					{typ: expWildcard, body: in[:pos]},
					{body: in[pos:]},
				}
			}
			return []expression{{typ: expWildcard, body: in}}
		} else {
			pos := asciiSet.LastIndex(in) + 1
			if pos < len(in) {
				return []expression{
					{body: in[:start]},
					{typ: expWildcard, body: in[start:pos]},
					{body: in[pos:]},
				}
			}
			return []expression{
				{body: in[:start]},
				{typ: expWildcard, body: in[start:]},
			}
		}
	}

	count := strings.Count(in, "[") + strings.Count(in, "{") + 2
	exps := make([]expression, 0, count)
	for {
		if start == -1 {
			exps = append(exps, expression{body: in})
			break
		} else if stop == -1 {
			if start > 0 {
				exps = append(exps, expression{body: in[:start]})
			}
			exps = append(exps, expression{typ: expWildcard, body: in[start:]})
			break
		} else if start > 0 {
			exps = append(exps, expression{body: in[:start]})
		}

		stop++
		e := getExpression(in[start:stop])
		if e.typ == expString {
			if len(exps) > 0 && exps[len(exps)-1].typ == expString {
				exps[len(exps)-1].body += e.body
			} else if e.body != "" {
				exps = append(exps, e)
			}
		} else {
			exps = append(exps, e)
		}

		in = in[stop:]
		if in == "" {
			break
		}
		start, stop = getPair(in)
	}

	last := len(exps) - 1
	if exps[last].typ == expWildcard {
		s := exps[last].body
		pos := asciiSet.LastIndex(s) + 1
		if pos < len(s) {
			exps[last].body = s[:pos]
			exps = append(exps, expression{body: s[pos:]})
		}
	}

	return exps
}

func expand(exps []expression, result []string, count, max, depth int, buf []byte) ([]string, []byte, error) {
	if len(exps) == 0 {
		// end of expressions, write result string
		result = append(result, string(buf))
		return result, buf, nil
	}

	var err error
	cur := len(buf)

	if max > 0 {
		if count <= 0 {
			count = 1
		}
		count *= exps[0].count()
		if max < count {
			max = 0
		}
	}
	if depth > 0 {
		if exps[0].count() > 1 {
			depth--
			if depth == 0 {
				max = 0
			}
		}
	}

	if max == 0 {
		buf = append(buf, exps[0].body...)
		result, buf, err = expand(exps[1:], result, count, max, depth, buf)
		if err != nil {
			return nil, buf, err
		}
	} else {
		for {
			buf, err = exps[0].appendNext(buf)
			if err == io.EOF {
				exps[0].reset()
				break
			} else if err != nil {
				return nil, buf, err
			}

			result, buf, err = expand(exps[1:], result, count, max, depth, buf)
			buf = buf[:cur]
			if err != nil {
				return nil, buf, err
			}
		}
	}

	return result, buf, nil
}
