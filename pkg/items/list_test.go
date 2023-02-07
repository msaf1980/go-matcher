package items

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func Test_listExpand(t *testing.T) {
	tests := []struct {
		s          string
		wantList   []string
		wantFailed bool
	}{
		{"{}", nil, false},
		{"{abc}", []string{"abc"}, false},
		{"{abc,z}", []string{"abc", "z"}, false},
		// partially broken
		{"{,abc}", []string{"abc"}, false},
		{"{abc,}", []string{"abc"}, false},
		{"{abc,,,q}", []string{"abc", "q"}, false},
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
			assert.Equal(t, tt.wantFailed, gotSuccess)
		})
	}
}
