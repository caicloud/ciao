package kubeflow

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/ciao/pkg/types"
)

func (b Backend) createPyTorchJob(parameter *types.Parameter) (*types.Job, error) {
	pytorchJob := b.Generator.GeneratePyTorchJob(parameter)
	pytorchJob, err := b.PyTorchJobClient.KubeflowV1alpha2().PyTorchJobs(metav1.NamespaceDefault).Create(pytorchJob)
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
