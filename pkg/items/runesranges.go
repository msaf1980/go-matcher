package items

import (
	"math/rand"
	"strings"

	"github.com/msaf1980/go-matcher/pkg/utils"
)

// RunesRanges is a range of rune symbols: [a-crzA-Z] is a range of a-c r z A-Z
type RunesRanges struct {
	utils.RunesRanges
}

func NewRunesRanges(ranges string) *RunesRanges {
	if rs, ok := utils.RunesRangeExpand(ranges); ok {
		return &RunesRanges{RunesRanges: rs}
	} else {
		return nil
	}
}

func (item *RunesRanges) Equal(a Item) bool {
	if v, ok := a.(*RunesRanges); ok {
		return item.RunesRanges.Equal(&v.RunesRanges)
	}
	return false
}

func (item *RunesRanges) WriteRandom(buf *strings.Builder) {
	if !item.ASCII.IsEmpty() {
		var i byte
		for i = 1; i < 128; i++ {
			if item.ASCII.Contains(i) {
				buf.WriteByte(i)
			}
		}
	}
	buf.WriteRune(item.UnicodeRanges[rand.Intn(len(item.UnicodeRanges))].First)
}

func (item *RunesRanges) WriteString(buf *strings.Builder) string {
	l := buf.Len()
	item.RunesRanges.WriteString(buf)
	return buf.String()[l:]
}

func (item *RunesRanges) String() string {
	var buf strings.Builder
	return item.WriteString(&buf)
}

func (item *RunesRanges) MinLen() int {
	return item.MinSize
}

func (item *RunesRanges) MaxLen() int {
	return item.MaxSize
}

func (item *RunesRanges) Find(s string) (index, length int, support FindFlag) {
	index, _, length = item.Index(s)
	return
}

func (item *RunesRanges) Match(s string) (offset int, support FindFlag) {
	_, offset = item.StartsWith(s)
	return
}

func (item *RunesRanges) MatchLast(s string) (offset int, support FindFlag) {
	_, offset = item.EndsWith(s)
	return
}
