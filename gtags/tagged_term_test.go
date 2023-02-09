package gtags

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

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
		require.NoError(b, err)
	}
}

func BenchmarkPathTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := PathTags(pathTags)
		require.NoError(b, err)
	}
}
