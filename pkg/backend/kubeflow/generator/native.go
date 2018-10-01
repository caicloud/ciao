package generator

import (
	pytorchv1alpha2 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1alpha2"
	tfv1alpha2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/ciao/pkg/types"
)

// Native is the type for native generator.
type Native struct {
}

// New returns a new native generator.
func NewNative() *Native {
	return &Native{}
}

// GenerateTFJob generates a new TFJob.
func (n Native) GenerateTFJob(parameter *types.Parameter) *tfv1alpha2.TFJob {
	psCount := int32(parameter.PSCount)
	workerCount := int32(parameter.WorkerCount)

	return &tfv1alpha2.TFJob{
		TypeMeta: metav1.TypeMeta{
			Kind: tfv1alpha2.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: tfv1alpha2.TFJobSpec{
			TFReplicaSpecs: map[tfv1alpha2.TFReplicaType]*tfv1alpha2.TFReplicaSpec{
				tfv1alpha2.TFReplicaTypePS: &tfv1alpha2.TFReplicaSpec{
					Replicas: &psCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:  defaultContainerNameTF,
									Image: parameter.Image,
								},
							},
						},
					},
				},
				tfv1alpha2.TFReplicaTypeWorker: &tfv1alpha2.TFReplicaSpec{
					Replicas: &workerCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:  defaultContainerNameTF,
									Image: parameter.Image,
								},
							},
						},
					},
				},
			},
		},
	}
}

// GeneratePyTorchJob generates a new PyTorchJob.
func (n Native) GeneratePyTorchJob(parameter *types.Parameter) *pytorchv1alpha2.PyTorchJob {
	masterCount := int32(parameter.MasterCount)
	workerCount := int32(parameter.WorkerCount)

	return &pytorchv1alpha2.PyTorchJob{
		TypeMeta: metav1.TypeMeta{
			Kind: pytorchv1alpha2.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: pytorchv1alpha2.PyTorchJobSpec{
			PyTorchReplicaSpecs: map[pytorchv1alpha2.PyTorchReplicaType]*pytorchv1alpha2.PyTorchReplicaSpec{
				pytorchv1alpha2.PyTorchReplicaTypeMaster: &pytorchv1alpha2.PyTorchReplicaSpec{
					Replicas: &masterCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:  defaultContainerNamePyTorch,
									Image: parameter.Image,
								},
							},
						},
					},
				},
				pytorchv1alpha2.PyTorchReplicaTypeWorker: &pytorchv1alpha2.PyTorchReplicaSpec{
					Replicas: &workerCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:  defaultContainerNamePyTorch,
									Image: parameter.Image,
								},
							},
						},
					},
				},
			},
		},
	}
}
