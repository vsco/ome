package pod

import (
	"context"
	"testing"

	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclientfake "sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/sgl-project/ome/pkg/apis/ome/v1beta1"
	"github.com/sgl-project/ome/pkg/constants"
)

func TestInjectFineTunedAdapter_SkipsWhenFineTunedWeightIsPVC(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	scheme := runtime.NewScheme()
	g.Expect(v1beta1.AddToScheme(scheme)).To(gomega.Succeed())

	pvcURI := "pvc://ome-models-efs/adapters/qwen3-vl-8b-lora"
	ftw := &v1beta1.FineTunedWeight{
		ObjectMeta: metav1.ObjectMeta{Name: "qwen3-vl-8b-lora"},
		Spec: v1beta1.FineTunedWeightSpec{
			Storage: &v1beta1.StorageSpec{
				StorageUri: &pvcURI,
			},
		},
	}
	cl := ctrlclientfake.NewClientBuilder().WithScheme(scheme).WithObjects(ftw).Build()

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "dummy", Namespace: constants.OMENamespace},
		Data:       map[string]string{fineTunedAdapterConfigMapKeyName: "{}"},
	}
	fa := newFineTunedAdapterInjector(cm, cl)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				constants.FineTunedAdapterInjectionKey: "qwen3-vl-8b-lora",
			},
		},
		Spec: v1.PodSpec{Containers: []v1.Container{{Name: constants.MainContainerName}}},
	}
	g.Expect(fa.InjectFineTunedAdapter(pod)).To(gomega.Succeed())
	g.Expect(pod.Spec.InitContainers).To(gomega.BeEmpty())
}

func TestInjectFineTunedAdapter_RequiresOCIConfigWhenNotPVC(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	scheme := runtime.NewScheme()
	g.Expect(v1beta1.AddToScheme(scheme)).To(gomega.Succeed())

	s3URI := "s3://bucket/prefix"
	ftw := &v1beta1.FineTunedWeight{
		ObjectMeta: metav1.ObjectMeta{Name: "ft-s3"},
		Spec: v1beta1.FineTunedWeightSpec{
			Storage: &v1beta1.StorageSpec{
				StorageUri: &s3URI,
			},
		},
	}
	cl := ctrlclientfake.NewClientBuilder().WithScheme(scheme).WithObjects(ftw).Build()

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "dummy", Namespace: constants.OMENamespace},
		Data:       map[string]string{fineTunedAdapterConfigMapKeyName: "{}"},
	}
	fa := newFineTunedAdapterInjector(cm, cl)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				constants.FineTunedAdapterInjectionKey: "ft-s3",
			},
		},
		Spec: v1.PodSpec{Containers: []v1.Container{{Name: constants.MainContainerName}}},
	}
	err := fa.InjectFineTunedAdapter(pod)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("FineTunedAdapterInjector"))
}
