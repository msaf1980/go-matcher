package gtags

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParsePath(t *testing.T) {
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
				t.Errorf("ParsePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTags, tt.wantTags) {
				t.Errorf("ParsePath() = %s", cmp.Diff(tt.wantTags, gotTags))
			}
		})
	}
}
