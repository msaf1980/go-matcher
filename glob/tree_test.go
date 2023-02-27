package glob

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

type verify struct {
	glob  string
	index int
}

func mergeVerify(globs []string, index []int) []verify {
	if len(globs) != len(index) {
		return nil
	}
	v := make([]verify, len(globs))
	for i := 0; i < len(globs); i++ {
		v[i].glob = globs[i]
		v[i].index = index[i]
	}
	return v
}

type testGlobTree struct {
	globs []string
	want  *GlobTree
	match map[string][]string
}

func runTestGlobTree(t *testing.T, n int, tt testGlobTree) {
	t.Run(fmt.Sprintf("%d#%#v", n, tt.globs), func(t *testing.T) {
		gtree := NewTree()
		for i, g := range tt.globs {
			_, _, err := gtree.AddGlob(g, i)

			if err != nil && err != ErrGlobExist {
				t.Fatalf("GlobTree.Add(%q) error = %v", g, err)
			}
		}

		if !reflect.DeepEqual(gtree, tt.want) {
			t.Fatalf("GlobTree(%#v) = %s", tt.globs, cmp.Diff(tt.want, gtree))
		}
		verifyGlobTree(t, tt.globs, tt.match, gtree)
	})
}

func verifyGlobTree(t *testing.T, globs []string, match map[string][]string, gtree *GlobTree) {
	for path, wantGlobs := range match {
		t.Run("#path="+path, func(t *testing.T) {
			var (
				globs []string
				index []int
			)
			first := -1
			matched := gtree.Match(path, &globs, &index, &first)

			verify := mergeVerify(globs, index)

			sort.Strings(globs)
			sort.Strings(wantGlobs)
			sort.Ints(index)

			if !reflect.DeepEqual(wantGlobs, globs) {
				t.Fatalf("GlobTree(%#v).Match(%q) globs = %s", globs, path, cmp.Diff(wantGlobs, globs))
			}

			if matched != len(globs) || len(globs) != len(index) {
				t.Fatalf("GlobTree(%#v).Match(%q) = %d, want %d, index = %d", globs, path, matched, len(globs), len(index))
			}

			for _, v := range verify {
				if v.glob != gtree.GlobsIndex[v.index] {
					t.Errorf("GlobTree(%#v).Match(%q) index = %d glob = %s, want %s",
						globs, path, v.index, gtree.GlobsIndex[v.index], v.glob)
				}
			}

			if len(index) > 0 {
				if first != index[0] {
					t.Errorf("GlobTree(%#v).Match(%q) first index = %d, want %d",
						globs, path, first, index[0])
				}
			}
		})
	}
}

func parseGlobs(globs []string) (g []*Glob) {
	for _, glob := range globs {
		g = append(g, ParseMust(glob))
	}
	return
}

func buildGlobSRegexp(globs []string) (re []*regexp.Regexp) {
	for _, glob := range globs {
		re = append(re, tests.BuildGlobRegexp(glob))
	}
	return
}
