package generator

import (
	pytorchv1alpha2 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1alpha2"
	tfv1alpha2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"

	"github.com/caicloud/ciao/pkg/types"
)

// Interface is the type for generator, which is used to generate TFJob/PyTorchJob.
type Interface interface {
	GenerateTFJob(parameter *types.Parameter) *tfv1alpha2.TFJob
	GeneratePyTorchJob(parameter *types.Parameter) *pytorchv1alpha2.PyTorchJob
}
