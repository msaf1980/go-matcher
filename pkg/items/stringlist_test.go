package items

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/msaf1980/go-matcher/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ ItemList = &StringList{}

func Test_listExpand(t *testing.T) {
	tests := []struct {
		s          string
		wantList   []string
		wantFailed bool
	}{
		{"{}", nil, false},
		{"{abc}", []string{"abc"}, false},
		{"{abc,z}", []string{"abc", "z"}, false},
		// has empty item
		{"{,abc}", []string{"", "abc"}, false},
		{"{abc,}", []string{"", "abc"}, false},
		{"{abc,,,q}", []string{"", "abc", "q"}, false},
		// duplicate
		{"{a,a}", []string{"a"}, false},
		{"{a,b,a}", []string{"a", "b"}, false},
		{"{b,a,b}", []string{"a", "b"}, false},
		{"{b,a,b,z}", []string{"a", "b", "z"}, false},
		{"{c,a,b,a,c,,z}", []string{"", "a", "b", "c", "z"}, false},
		// broken
		{"", nil, true},
		{"{a,", nil, true},
		{"a}", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			gotList, gotSuccess := ListExpand(tt.s)
			if !reflect.DeepEqual(gotList, tt.wantList) {
				t.Errorf("listExpand(%s) = %s", tt.s, cmp.Diff(tt.wantList, gotList))
			}
			assert.Equal(t, !tt.wantFailed, gotSuccess)
		})
	}
}

type matchN struct {
	s      string
	n      int
	index  int
	length int
}

func TestStringList(t *testing.T) {
	tests := []struct {
		list          string
		wantItem      Item
		wantFailed    bool
		wantMinLen    int
		wantMaxLen    int
		wantMatchFlag FindFlag
		wantFindFlag  FindFlag
		// List
		wantCanEmpty   bool
		wantMatchFirst map[string]bool
		wantMatchN     []matchN
		wantFindFirst  map[string]int
		wantFindN      []matchN
	}{
		{
			list: "{f,bc}",
			wantItem: &StringList{
				Vals: []string{"bc", "f"}, MinSize: 1, MaxSize: 2,
				ASCIIStarted: true, FirstASCII: utils.MakeASCIISetMust("fb"),
			},
			wantMinLen: 1, wantMaxLen: 2,
			wantMatchFlag: FindList,
			wantFindFlag:  FindList,
			wantMatchFirst: map[string]bool{
				"find": true, "b": true, "c": false, "cb": false, "cf": false,
				"ФЫВz": false, // unicode match
			},
			wantMatchN: []matchN{
				{s: "", n: 0, index: -1},
				{s: "fbci", n: 0, index: -1},
				{s: "bcfi", n: 0, index: 2},
				{s: "fbci", n: 1, index: 1},
			},
			wantFindFirst: map[string]int{
				"find": 0, "b": 0, "c": -1, "cb": 1, "cf": 1, "Яb": 2, "de": -1, "Яd": -1,
			},
			wantFindN: []matchN{
				{s: "", n: 0, index: -1},
				{s: "fbci", n: 0, index: 1, length: 2},
				{s: "fbci", n: 1, index: 0, length: 1},
			},
		},
		{
			list: "{f,Яc}",
			wantItem: &StringList{
				Vals: []string{"f", "Яc"}, MinSize: 1, MaxSize: 3,
			},
			wantMinLen: 1, wantMaxLen: 3,
			wantMatchFlag: FindList,
			wantFindFlag:  FindList,
			wantMatchFirst: map[string]bool{
				"find": true, "b": true, "c": true, "cb": true, "cf": true, "ЯЫВz": true, // match skipped
			},
			wantMatchN: []matchN{
				{s: "", n: 0, index: -1},
				{s: "fbci", n: 0, index: 1},
				{s: "Яcfi", n: 0, index: -1},
				{s: "Яcfbci", n: 1, index: 3},
			},
			wantFindFirst: map[string]int{
				"find": 0, "b": 0, "c": 0, "cb": 0, "cf": 0, "Яb": 0, "de": 0, "Яd": 0,
			},
			wantFindN: []matchN{
				{s: "", n: 0, index: -1},
				{s: "fbci", n: 0, index: 0, length: 1},
				{s: "fbci", n: 1, index: -1},
				{s: "fЯcbci", n: 1, index: 1, length: 3},
			},
		},
		{
			list: "{f,bЯ}",
			wantItem: &StringList{
				Vals: []string{"bЯ", "f"}, MinSize: 1, MaxSize: 3,
				ASCIIStarted: true, FirstASCII: utils.MakeASCIISetMust("fb"),
			},
			wantMinLen: 1, wantMaxLen: 3,
			wantMatchFlag: FindList,
			wantFindFlag:  FindList,
			wantMatchFirst: map[string]bool{
				"find": true, "b": true, "bc": true, "bЯЫВz": true,
				"c": false, "cb": false, "cf": false, "ЯЫВz": false,
			},
			wantMatchN: []matchN{
				{s: "", n: 0, index: -1},
				{s: "Яcfbci", n: 0, index: -1},
				{s: "bЯcfbci", n: 0, index: 3},
				{s: "Яcfi", n: 1, index: -1},
				{s: "fbci", n: 1, index: 1},
			},
			wantFindFirst: map[string]int{
				"find": 0, "b": 0, "bD": 0, "cb": 1, "cf": 1, "bЯ": 0, "Яb": 2, "c": -1, "de": -1, "Яd": -1,
			},
			wantFindN: []matchN{
				{s: "", n: 0, index: -1},
				{s: "fbci", n: 0, index: -1},
				{s: "fЯcbci", n: 0, index: -1},
				{s: "fbЯbci", n: 0, index: 1, length: 3},
				{s: "fbci", n: 1, index: 0, length: 1},
				{s: "bЯfbci", n: 1, index: 3, length: 1},
			},
		},
		{list: "{f}", wantItem: Byte('f'), wantMinLen: 1, wantMaxLen: 1},
		{list: "{Я}", wantItem: Rune('Я'), wantMinLen: 2, wantMaxLen: 2},
		{list: "{Яf}", wantItem: NewString("Яf"), wantMinLen: 3, wantMaxLen: 3},
		{list: "", wantFailed: true}, // broken
		{list: "{}"},                 // empty return nil
	}
	for n, tt := range tests {
		t.Run(strconv.Itoa(n)+"#"+tt.list, func(t *testing.T) {
			list, ok := ListExpand(tt.list)
			require.Equal(t, !tt.wantFailed, ok, "ListExpand(%q) status", tt.list)

			item, _ := NewItemList(list)
			if !reflect.DeepEqual(tt.wantItem, item) {
				t.Fatalf("NewItemList(%q) = %s", tt.list, cmp.Diff(tt.wantItem, item))
			}
			if item != nil {
				if tt.wantMinLen != item.MinLen() {
					t.Errorf("NewItemList(%q) minLen = %d, want %d", tt.list, item.MinLen(), tt.wantMinLen)
				}
				if tt.wantMaxLen != item.MaxLen() {
					t.Errorf("NewItemList(%q) maxLen = %d, want %d", tt.list, item.MaxLen(), tt.wantMaxLen)
				}

				index, flag := item.Match("")
				if flag != tt.wantMatchFlag {
					t.Errorf("List(%q).Match(%q) flag = %v, want %v", tt.list, "", flag, tt.wantMatchFlag)
				}
				if flag == FindList {
					if index != 0 {
						t.Errorf("List(%q).Match(%q) = %v, want %v", tt.list, "", index, 0)
					}
					l := item.(ItemList)
					assert.Equal(t, tt.wantCanEmpty, l.IsOptional(), "List(%q) can empty", tt.list)
					if len(tt.wantMatchFirst) > 0 {
						for s, wantOK := range tt.wantMatchFirst {
							if ok, _ := l.MatchFirst(s); ok != wantOK {
								t.Errorf("List(%s).MatchFirst(%q) = %v, want %v", l.String(), s, ok, wantOK)
							}
						}
					}
					for _, m := range tt.wantMatchN {
						if offset := l.MatchN(m.s, m.n); offset != m.index {
							t.Errorf("List(%s).MatchN(%q, %d) = %d, want %d", l.String(), m.s, m.n, offset, m.index)
						}
					}

					if len(tt.wantFindFirst) > 0 {
						for s, wantIndex := range tt.wantFindFirst {
							if index, _ := l.FindFirst(s); index != wantIndex {
								t.Errorf("List(%s).FindFirst(%q) = %v, want %v", l.String(), s, index, wantIndex)
							}
						}
					}
					for _, m := range tt.wantFindN {
						index, length := l.FindN(m.s, m.n)
						if index != m.index {
							t.Errorf("List(%s).FindN(%q, %d) index = %d, want %d", l.String(), m.s, m.n, index, m.index)
						}
						if length != m.length {
							t.Errorf("List(%s).FindN(%q, %d) length = %d, want %d", l.String(), m.s, m.n, length, m.length)
						}
					}
				}

				index, length, flag := item.Find("")
				if flag != tt.wantMatchFlag {
					t.Errorf("List(%q).Find(%q) flag = %v, want %v", tt.list, "", flag, tt.wantMatchFlag)
				}
				if flag == FindList {
					if index != 0 {
						t.Errorf("List(%q).Match(%q) = %v, want %v", tt.list, "", index, 0)
					}
					if length != 0 {
						t.Errorf("List(%q).Match(%q) length = %v, want %v", tt.list, "", length, 0)
					}
				}
			}
		})
	}
}

var (
	listStringListASCII   = "{ac,c,b}"
	stringStringListASCII = strings.Repeat("find", 60) + "bac"
)

func BenchmarkStringList_Find_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list, _ := ListExpand(listStringListASCII)
		if list == nil {
			b.Fail()
		}

		item, _ := NewItemList(list)
		l := item.(ItemList)

		s := stringStringListASCII
		index, _ := l.FindFirst(s)
		if index == -1 {
			b.Fatal("index")
		}
		s = s[index:]
		for i := 0; i < l.Len(); i++ {
			if offset, _ := l.FindN(s, i); offset == -1 {
				b.Fatal("offset")
			}
		}
	}
}

func BenchmarkStringList_Find_ASCII_Prealloc(b *testing.B) {
	list, _ := ListExpand(listStringListASCII)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListASCII
		for i := 0; i < l.Len(); i++ {
			if offset, _ := l.FindN(s, i); offset == -1 {
				b.Fatal("offset")
			}
		}
	}
}

func BenchmarkStringList_Find_ASCII_Skip(b *testing.B) {
	list, _ := ListExpand(listStringListASCII)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListASCII
		index, _ := l.FindFirst(s)
		if index == -1 {
			b.Fatal("index")
		}
		s = s[index:]
		for i := 0; i < l.Len(); i++ {
			if offset, _ := l.FindN(s, i); offset == -1 {
				b.Fatal("offset")
			}
		}
	}
}

func BenchmarkStringList_Match_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list, _ := ListExpand(listStringListASCII)
		if list == nil {
			b.Fail()
		}

		item, _ := NewItemList(list)
		l := item.(ItemList)

		s := stringStringListASCII
		_, _ = l.MatchFirst(s)

		for i := 0; i < l.Len(); i++ {
			_ = l.MatchN(s, i)
		}

	}
}

func BenchmarkStringList_Match_ASCII_Prealloc(b *testing.B) {
	list, _ := ListExpand(listStringListASCII)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListASCII
		_, _ = l.MatchFirst(s)
		for i := 0; i < l.Len(); i++ {
			_ = l.MatchN(s, i)
		}
	}
}

func BenchmarkStringList_Match_ASCII_Skip(b *testing.B) {
	list, _ := ListExpand(listStringListASCII)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListASCII
		ok, sup := l.MatchFirst(s)
		if !ok || !sup {
			for i := 0; i < l.Len(); i++ {
				_ = l.MatchN(s, i)
			}
		}
	}
}

var (
	listStringListUnicode   = "{ac,c,b,Я}"
	stringStringListUnicode = strings.Repeat("find", 60) + "Яbac"
)

func BenchmarkStringList_Find_Unicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list, _ := ListExpand(listStringListUnicode)
		if list == nil {
			b.Fail()
		}

		item, _ := NewItemList(list)
		l := item.(ItemList)

		s := stringStringListUnicode
		index, _ := l.FindFirst(s)
		if index == -1 {
			b.Fatal("index")
		}
		s = s[index:]
		for i := 0; i < l.Len(); i++ {
			if offset, _ := l.FindN(s, i); offset == -1 {
				b.Fatal("offset")
			}
		}
	}
}

func BenchmarkStringList_Find_Unicode_Prealloc(b *testing.B) {
	list, _ := ListExpand(listStringListUnicode)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListUnicode
		for i := 0; i < l.Len(); i++ {
			if offset, _ := l.FindN(s, i); offset == -1 {
				b.Fatal("offset")
			}
		}
	}
}

func BenchmarkStringList_Find_Unicode_Skip(b *testing.B) {
	list, _ := ListExpand(listStringListUnicode)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListUnicode
		index, _ := l.FindFirst(s)
		if index == -1 {
			b.Fatal("index")
		}
		s = s[index:]
		for i := 0; i < l.Len(); i++ {
			if offset, _ := l.FindN(s, i); offset == -1 {
				b.Fatal("offset")
			}
		}
	}
}

func BenchmarkStringList_Match_Unicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list, _ := ListExpand(listStringListUnicode)
		if list == nil {
			b.Fail()
		}

		item, _ := NewItemList(list)
		l := item.(ItemList)

		s := stringStringListUnicode
		_, _ = l.MatchFirst(s)

		for i := 0; i < l.Len(); i++ {
			_ = l.MatchN(s, i)
		}
	}
}

func BenchmarkStringList_Match_Unicode_Prealloc(b *testing.B) {
	list, _ := ListExpand(listStringListUnicode)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListUnicode
		_, _ = l.MatchFirst(s)

		for i := 0; i < l.Len(); i++ {
			_ = l.MatchN(s, i)
		}
	}
}

func BenchmarkStringList_Match_Unicode_Skip(b *testing.B) {
	list, _ := ListExpand(listStringListUnicode)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListUnicode
		ok, sup := l.MatchFirst(s)
		if !ok || !sup {
			for i := 0; i < l.Len(); i++ {
				_ = l.MatchN(s, i)
			}
		}
	}
}

var (
	stringStringListMiss = strings.Repeat("find", 60)
)

func BenchmarkStringList_Match_ASCII_Miss(b *testing.B) {
	list, _ := ListExpand(listStringListASCII)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListMiss
		for i := 0; i < l.Len(); i++ {
			_ = l.MatchN(s, i)
		}
	}
}

func BenchmarkStringList_Match_ASCII_Miss_Skip(b *testing.B) {
	list, _ := ListExpand(listStringListASCII)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListMiss
		ok, sup := l.MatchFirst(s)
		if !ok || !sup {
			for i := 0; i < l.Len(); i++ {
				_ = l.MatchN(s, i)
			}
		}
	}
}

func BenchmarkStringList_Match_Unicode_Miss(b *testing.B) {
	list, _ := ListExpand(listStringListUnicode)
	if list == nil {
		b.Fail()
	}

	item, _ := NewItemList(list)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := item.(ItemList)

		s := stringStringListMiss
		ok, sup := l.MatchFirst(s)
		if !ok || !sup {
			for i := 0; i < l.Len(); i++ {
				_ = l.MatchN(s, i)
			}
		}
	}
}
