package glob

import (
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type testGlob struct {
	glob    string
	want    *Glob
	verify  string
	wantErr bool
	match   []string
	miss    []string
}

func runTestGlob(t *testing.T, n int, tt testGlob) {
	t.Run(strconv.Itoa(n)+"#"+tt.glob, func(t *testing.T) {
		g, err := Parse(tt.glob)

		if (err != nil) != tt.wantErr {
			t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.glob, err, tt.wantErr)
		}
		if tt.wantErr {
			assert.Equal(t, 0, len(tt.match), "can't check on error")
			assert.Equal(t, 0, len(tt.miss), "can't check on error")
		} else {
			if !reflect.DeepEqual(g, tt.want) {
				t.Fatalf("Glob(%q) = %s", tt.glob, cmp.Diff(tt.want, g))
			}
			verifyGlob(t, tt.match, tt.miss, g, tt.verify)
		}
	})
}

func verifyGlob(t *testing.T, match []string, miss []string, g *Glob, verifyRegexp string) {
	var re *regexp.Regexp
	if verifyRegexp != "" {
		re = regexp.MustCompile(verifyRegexp)
	}
	for n, path := range match {
		t.Run(strconv.Itoa(n)+"#path="+path, func(t *testing.T) {
			matched := g.Match(path)
			if re == nil {
				if !matched {
					t.Errorf("Glob(%q).Match(%q) = %v, want true", g.Node, path, matched)
				}
			} else {
				reMatched := re.MatchString(path)
				if !matched || matched != reMatched {
					t.Errorf("Glob(%q).Match(%q) = %v, want true, re(%q) = %v", g.Node, path, matched, verifyRegexp, reMatched)
				}
			}
		})
	}
	for n, path := range miss {
		t.Run(strconv.Itoa(n)+"#MISS#path="+path, func(t *testing.T) {
			matched := g.Match(path)
			if re == nil {
				if matched {
					t.Errorf("Glob(%q).Match(%q) = %v, want false", g.Node, path, matched)
				}
			} else {
				reMatched := re.MatchString(path)
				if matched || matched != reMatched {
					t.Errorf("Glob(%q).Match(%q) = %v, want false, re(%q) = %v", g.Node, path, matched, verifyRegexp, reMatched)
				}
			}
		})
	}
}
