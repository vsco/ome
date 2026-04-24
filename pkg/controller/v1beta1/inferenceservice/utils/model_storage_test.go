package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sgl-project/ome/pkg/apis/ome/v1beta1"
	"github.com/sgl-project/ome/pkg/constants"
	"github.com/sgl-project/ome/pkg/controller/v1beta1/controllerconfig"
)

func TestModelVolumeMountSubPathForPVC(t *testing.T) {
	root := constants.DefaultModelPVCMountRoot
	nested := filepath.Join(root, "a", "b")
	tests := []struct {
		name        string
		mountRoot   string
		storagePath string
		want        string
	}{
		{"equal root", root, root, ""},
		{"subdir", root, filepath.Join(root, "qwen3-vl-8b-instruct"), "qwen3-vl-8b-instruct"},
		{"nested", root, nested, filepath.Join("a", "b")},
		{"outside tree", root, "/other/path", ""},
		{"empty root", "", filepath.Join(root, "x"), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ModelVolumeMountSubPathForPVC(tt.mountRoot, tt.storagePath)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestModelStoragePVCClaimName(t *testing.T) {
	cfg := &controllerconfig.InferenceServicesConfig{
		ModelStorage: controllerconfig.ModelStorageConfig{PVCClaimName: "from-config"},
	}
	isvc := &v1beta1.InferenceService{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{constants.ModelStoragePVCAnnotationKey: "  from-annotation  "},
		},
	}
	assert.Equal(t, "from-annotation", ModelStoragePVCClaimName(isvc, cfg))

	isvc2 := &v1beta1.InferenceService{ObjectMeta: metav1.ObjectMeta{}}
	assert.Equal(t, "from-config", ModelStoragePVCClaimName(isvc2, cfg))
	assert.Equal(t, "", ModelStoragePVCClaimName(isvc2, nil))
}
