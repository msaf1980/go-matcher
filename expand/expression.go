package expand

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

type expTyp int8

const (
	expString expTyp = iota
	expList
	expRunes
)

var (
	expTypMap = map[int]string{
		0: "",
		1: "{",
		2: "[",
	}
)

type ErrUnsupportedExp expTyp

func (e ErrUnsupportedExp) Error() string {
	t := int(e)
	n, ok := expTypMap[t]
	if !ok {
		n = strconv.Itoa(t)
	}
	return "unsupported expression expand: " + n
}

// expression represents all possible expandable types
type expression struct {
	typ   expTyp
	body  string
	list  []string
	runes []runes
	pos   int // -1 for EOF
}

func (e *expression) count() int {
	n := 0
	switch e.typ {
	case expString:
		n = 1
	case expList:
		n = len(e.list)
	case expRunes:
		for i := 0; i < len(e.runes); i++ {
			n += e.runes[i].count()
		}
	default:
		panic(fmt.Errorf("BUG: not implemented for %d", e.typ))
	}
	return n
}

func (e *expression) reset() {
	e.pos = 0
	for i := 0; i < len(e.runes); i++ {
		e.runes[i].reset()
	}
}

func (e *expression) appendNext(out []byte) ([]byte, error) {
	if e.pos == -1 {
		return out, io.EOF
	}
	switch e.typ {
	case expString:
		out = append(out, e.body...)
		e.pos = -1
	case expList:
		out = append(out, e.list[e.pos]...)
		if e.pos < len(e.list)-1 {
			e.pos++
		} else {
			e.pos = -1
		}
	case expRunes:
		for {
			if c := e.runes[e.pos].next(); c == utf8.RuneError {
				if e.pos < len(e.runes)-1 {
					e.pos++
				} else {
					return out, io.EOF
				}
				e.runes[e.pos].reset()
			} else {
				out = utf8.AppendRune(out, c)
				break
			}
		}
	default:
		panic(fmt.Errorf("BUG: not implemented for %d", e.typ))
	}
	return out, nil
}

// getExpression returns expression depends on the input
func getExpression(in string) expression {
	orig := in
	in = in[1 : len(in)-1]
	if len(in) == 0 {
		return expression{body: in}
	}
	switch orig[0] {
	case '{':
		if strings.ContainsRune(in, ',') {
			return expression{typ: expList, body: orig, list: strings.Split(in, ",")}
		} else {
			return expression{body: in}
		}
	case '[':
		if len(in) == 1 {
			return expression{body: in}
		} else {
			// TODO
			// return rune{in}
			rs, ok := runesRangeExpand(in)
			if !ok {
				return expression{body: in}
			}
			return expression{typ: expRunes, body: orig, runes: rs}
		}
	default:
		return expression{body: orig}
	}

	// args := strings.Split(in, dots)
	// if len(args) != 2 && len(args) != 3 {
	// 	return none{orig}
	// }

	// if len(args) != 2 {
	// 	return none{orig}
	// }
	// // rSeq := make([]rune, len(args))
	// // for i, a := range args {
	// // 	r := []rune(a)
	// // 	if len(r) != 1 {
	// // 		return none{orig}
	// // 	}
	// // 	rSeq[i] = r[0]
	// // }

	// // return runes{rSeq}
	// return none{orig}
}
