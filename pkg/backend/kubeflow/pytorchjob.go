package kubeflow

import (
	pytorchv1alpha2 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1alpha2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	defaultContainerNamePyTorch = "pytorch"
)

func (b Backend) createPyTorchJob(parameter *types.Parameter) (*types.Job, error) {
	pytorchJob := generatePyTorchJob(parameter)
	pytorchJob, err := b.PyTorchJobClient.KubeflowV1alpha2().PyTorchJobs(namespaceDefault).Create(pytorchJob)
	if err != nil {
		return nil, err
	}
	return &types.Job{
		Name:      pytorchJob.Name,
		Framework: types.FrameworkTypePyTorch,
		Master:    parameter.MasterCount,
		Worker:    parameter.WorkerCount,
	}, nil
}

func generatePyTorchJob(parameter *types.Parameter) *pytorchv1alpha2.PyTorchJob {
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
