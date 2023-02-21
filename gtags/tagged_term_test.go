package gtags

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type testTaggedTermList struct {
	query      string
	want       TaggedTermList
	wantQuery  string
	wantErr    bool
	matchPaths []string
	missPaths  []string
}

func runTestTaggedTermList(t *testing.T, tt testTaggedTermList) {
	terms, err := ParseSeriesByTag(tt.query)
	if (err != nil) != tt.wantErr {
		t.Fatalf("ParseSeriesByTag(%q) error = %v, wantErr %v", tt.query, err, tt.wantErr)
	}
	if err = terms.Build(); err == nil {
		if tt.query != "" {
			var buf strings.Builder
			buf.Grow(len(tt.query))
			terms.Rewrite(&buf)
			newQuery := buf.String()
			assert.Equal(t, tt.wantQuery, newQuery, tt.query)
		}

		if !cmp.Equal(terms, tt.want, cmpTransform) {
			t.Errorf("TagsMatcher.Add() = %s", cmp.Diff(tt.want, terms, cmpTransform))
		}
		verifyTaggedTermList(t, tt.matchPaths, tt.missPaths, terms)
	}
	if tt.wantErr {
		assert.Equal(t, 0, len(tt.matchPaths), "can't check on error", tt.query)
		assert.Equal(t, 0, len(tt.missPaths), "can't check on error", tt.query)
	}
}

func verifyTaggedTermList(t *testing.T, matchPaths, missPaths []string, terms TaggedTermList) {
	for _, path := range matchPaths {
		tags, err := PathTagsMap(path)
		if err != nil {
			t.Errorf("PathTags(%q) err = %q", path, err.Error())
		}
		if !terms.MatchByTagsMap(tags) {
			t.Errorf("TaggedTermList.MatchByPathMap(%q) != true", path)
		}
		tagsMap, err := PathTagsMap(path)
		if err != nil {
			t.Errorf("PathTagsMap(%q) err = %q", path, err.Error())
		}
		if !terms.MatchByTagsMap(tagsMap) {
			t.Errorf("TaggedTermList.MatchByPathMap(%q) != true", path)
		}
	}
	for _, path := range missPaths {
		tags, err := PathTags(path)
		if err != nil {
			t.Errorf("PathTags(%q) err = %q", path, err.Error())
		}
		if terms.MatchByTags(tags) {
			t.Errorf("TaggedTermList.MatchByPath(%q) != false", path)
		}
		tagsMap, err := PathTagsMap(path)
		if err != nil {
			t.Errorf("PathTagsMap(%q) err = %q", path, err.Error())
		}
		if terms.MatchByTagsMap(tagsMap) {
			t.Errorf("TaggedTermList.MatchByPathMap(%q) != false", path)
		}
	}
}

func TestPathTagsMap(t *testing.T) {
	tests := []struct {
		path     string
		wantTags map[string]string
		wantErr  bool
	}{
		{
			path: "kube_pod_status_phase?app_kubernetes_io_component=metrics&app_kubernetes_io_name=kube-state-metrics&app_kubernetes_io_part_of=kube-state-metrics&app_kubernetes_io_version=2.7.0&helm_sh_chart=kube-state-metrics-4.24.0&instance=192.168.0.85%3A8080&job=kubernetes-service-endpoints",
			wantTags: map[string]string{
				"__name__":                    "kube_pod_status_phase",
				"app_kubernetes_io_component": "metrics",
				"app_kubernetes_io_name":      "kube-state-metrics",
				"app_kubernetes_io_part_of":   "kube-state-metrics",
				"app_kubernetes_io_version":   "2.7.0",
				"helm_sh_chart":               "kube-state-metrics-4.24.0",
				"instance":                    "192.168.0.85:8080",
				"job":                         "kubernetes-service-endpoints",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			gotTags, err := PathTagsMap(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PathTagsMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTags, tt.wantTags) {
				t.Errorf("PathTagsMap() = %s", cmp.Diff(tt.wantTags, gotTags))
			}
		})
	}
}

func TestPathTags(t *testing.T) {
	tests := []struct {
		path     string
		wantTags []Tag
		wantErr  bool
	}{
		{
			path: "kube_pod_status_phase?app_kubernetes_io_component=metrics&app_kubernetes_io_name=kube-state-metrics&app_kubernetes_io_part_of=kube-state-metrics&app_kubernetes_io_version=2.7.0&helm_sh_chart=kube-state-metrics-4.24.0&instance=192.168.0.85%3A8080&job=kubernetes-service-endpoints",
			wantTags: []Tag{
				{"__name__", "kube_pod_status_phase"},
				{"app_kubernetes_io_component", "metrics"},
				{"app_kubernetes_io_name", "kube-state-metrics"},
				{"app_kubernetes_io_part_of", "kube-state-metrics"},
				{"app_kubernetes_io_version", "2.7.0"},
				{"helm_sh_chart", "kube-state-metrics-4.24.0"},
				{"instance", "192.168.0.85:8080"},
				{"job", "kubernetes-service-endpoints"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			gotTags, err := PathTags(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PathTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTags, tt.wantTags) {
				t.Errorf("PathTags() = %s", cmp.Diff(tt.wantTags, gotTags))
			}
		})
	}
}

var (
	pathTags = "kube_pod_status_phase?app_kubernetes_io_component=metrics&app_kubernetes_io_name=kube-state-metrics&app_kubernetes_io_part_of=kube-state-metrics&app_kubernetes_io_version=2.7.0&helm_sh_chart=kube-state-metrics-4.24.0&instance=192.168.0.85%3A8080&job=kubernetes-service-endpoints"
)

func BenchmarkPathTagsMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := PathTagsMap(pathTags)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPathTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := PathTags(pathTags)
		if err != nil {
			b.Fatal(err)
		}
	}
}
