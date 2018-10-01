package kubeflow

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/ciao/pkg/types"
)

func (b Backend) createTFJob(parameter *types.Parameter) (*types.Job, error) {
	tfJob := b.Generator.GenerateTFJob(parameter)
	tfJob, err := b.TFJobClient.KubeflowV1alpha2().TFJobs(metav1.NamespaceDefault).Create(tfJob)
	if err != nil {
		return nil, err
	}
	return &types.Job{
		Name:      tfJob.Name,
		Framework: types.FrameworkTypeTensorFlow,
		PS:        parameter.PSCount,
		Worker:    parameter.WorkerCount,
	}, nil
}
