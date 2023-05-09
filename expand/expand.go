package expand

import (
	"errors"
	"io"
	"strings"
)

var (
	ErrNotClosed = errors.New("expression not closed")
)

// Expand takes the string contains the shell expansion expression and returns list of strings after they are expanded.
//
// Argument max for restict max expanded results, > 0 - restuct  max expamnded results, 0 - disables expand, -1 - unlimited, -2 - expand only first node, -3 - expand only two nodes, etc.
func Expand(in string, max int) ([]string, error) {
	if max == 0 {
		return []string{in}, nil
	}
	start, stop := getPair(in)
	if start == -1 {
		return []string{in}, nil
	}

	count := (strings.Count(in, "[") + strings.Count(in, "{"))

	exps := make([]expression, 0, count)

	buf := make([]byte, 0, len(in))
	buf = append(buf, in[:start]...)
	for {
		in = in[start:]
		stop -= start
		exps = append(exps, getExpression(in[:stop+1]))
		in = in[stop+1:]
		if in == "" {
			break
		}
		start, stop = getPair(in)
		if start == -1 {
			exps = append(exps, expression{body: in})
			break
		} else {
			exps = append(exps, expression{body: in[:start]})
		}
	}

	count = 1
	for i := 0; i < len(exps); i++ {
		count *= exps[i].count()
	}
	// TODO: may be restrict max length
	result := make([]string, 0, count)

	var err error
	if result, _, err = expand(exps, result, 0, max, buf); err != nil {
		return nil, err
	}

	return result, nil
}

func expand(exps []expression, result []string, count, max int, buf []byte) ([]string, []byte, error) {
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

	if max == 0 {
		buf = append(buf, exps[0].body...)
		result, buf, err = expand(exps[1:], result, count, max, buf)
		if err != nil {
			return nil, buf, err
		}
	} else {
		if max == -2 {
			max = 0
		} else if max < -2 {
			max++
		}
		for {
			buf, err = exps[0].appendNext(buf)
			if err == io.EOF {
				exps[0].reset()
				break
			} else if err != nil {
				return nil, buf, err
			}

			result, buf, err = expand(exps[1:], result, count, max, buf)
			buf = buf[:cur]
			if err != nil {
				return nil, buf, err
			}
		}
	}

	return result, buf, nil
}

// getPair returns the top level expression. If the first `{` doesn't have the pair,
// it recursively executed for the substring after it
func getPair(in string) (start, stop int) {
	start = -1
	stop = -1
	for i, c := range in {
		switch c {
		case '{', '[':
			if start == -1 {
				start = i
			} else {
				start = -1
				return
			}
		case '}':
			if start == -1 || in[start] == '[' {
				start = -1
				return
			}
			stop = i
			return
		case ']':
			if start == -1 || in[start] == '{' {
				start = -1
				return
			}
			stop = i
			return
		}
	}

	return -1, -1
}
