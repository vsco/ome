package utils

import (
	"path/filepath"
	"strings"

	"github.com/sgl-project/ome/pkg/apis/ome/v1beta1"
	"github.com/sgl-project/ome/pkg/constants"
	"github.com/sgl-project/ome/pkg/controller/v1beta1/controllerconfig"
)

// ModelStoragePVCClaimName returns the PVC claim name for base-model storage when using a
// shared filesystem (e.g. EFS ReadWriteMany). Order: InferenceService annotation
// (ome.io/model-storage-pvc), then cluster ConfigMap modelStorage.pvcClaimName, else empty
// (hostPath mode).
func ModelStoragePVCClaimName(isvc *v1beta1.InferenceService, cfg *controllerconfig.InferenceServicesConfig) string {
	if isvc != nil {
		if v := strings.TrimSpace(isvc.Annotations[constants.ModelStoragePVCAnnotationKey]); v != "" {
			return v
		}
	}
	if cfg != nil {
		return strings.TrimSpace(cfg.ModelStorage.PVCClaimName)
	}
	return ""
}

// ModelStoragePVCMountRoot returns the directory inside the pod that corresponds to the root of
// the shared PVC (used with SubPath so per-model paths align with hostPath layout).
func ModelStoragePVCMountRoot(cfg *controllerconfig.InferenceServicesConfig) string {
	if cfg != nil && strings.TrimSpace(cfg.ModelStorage.PVCMountRoot) != "" {
		return filepath.Clean(cfg.ModelStorage.PVCMountRoot)
	}
	return constants.DefaultModelPVCMountRoot
}

// ModelVolumeMountSubPathForPVC returns the SubPath for a volumeMount when the model directory
// (storagePath) is under pvcMountRoot on the shared PVC. Empty means mount the whole claim at
// storagePath (single-directory layout).
func ModelVolumeMountSubPathForPVC(pvcMountRoot, storagePath string) string {
	if pvcMountRoot == "" || storagePath == "" {
		return ""
	}
	root := filepath.Clean(pvcMountRoot)
	path := filepath.Clean(storagePath)
	if root == path {
		return ""
	}
	sep := string(filepath.Separator)
	if path != root && !strings.HasPrefix(path, root+sep) {
		return ""
	}
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return ""
	}
	return rel
}
