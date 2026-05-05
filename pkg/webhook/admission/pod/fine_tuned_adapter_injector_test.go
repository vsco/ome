package pod

import (
	"testing"
)

func TestSkipFineTunedAdapterInit(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		uri  string
		skip bool
	}{
		{"empty", "", true},
		{"whitespace", "   ", true},
		{"pvc", "pvc://harmony-ome-models/qwen3-vl-8b-instruct", true},
		{"s3", "s3://vsco-lora-adapters-dev/studio/harmony/qwen3-vl-8b-instruct/v0p2", true},
		{"hf", "hf://org/repo@main", true},
		{"oci", "oci://namespace/bucket@object", false},
		{"oci_trimmed", "  oci://ns/bucket@obj  ", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := skipFineTunedAdapterInit(tc.uri); got != tc.skip {
				t.Fatalf("skipFineTunedAdapterInit(%q) = %v, want %v", tc.uri, got, tc.skip)
			}
		})
	}
}
