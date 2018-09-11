package kubeflow

import (
	"fmt"

	"k8s.io/api/core/v1"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	labelTFJobName      = "tf_job_name"
	tfReplicaTypeLabel  = "tf-replica-type"
	tfReplicaIndexLabel = "tf-replica-index"

	labelPyTorchJobName      = "pytorch_job_name"
	pytorchReplicaTypeLabel  = "pytorch-replica-type"
	pytorchReplicaIndexLabel = "pytorch-replica-index"
)

// GetLabelSelectorForJob gets label selector for the given job.
func GetLabelSelectorForJob(job *types.Job) string {
	switch job.Framework {
	case types.FrameworkTypeTensorFlow:
		return fmt.Sprintf("%s=%s", labelTFJobName, job.Name)
	case types.FrameworkTypePyTorch:
		return fmt.Sprintf("%s=%s", labelPyTorchJobName, job.Name)
	default:
		return "-1"
	}
}

// GetReplicaInstanceForPod gets the instance name of the given pod.
// e.g. kubeflow-xsadd-worker-0 will return worker-0.
func GetReplicaInstanceForPod(job *types.Job, pod v1.Pod) string {
	switch job.Framework {
	case types.FrameworkTypeTensorFlow:
		return fmt.Sprintf("%s-%s", pod.Labels[tfReplicaTypeLabel], pod.Labels[tfReplicaIndexLabel])
	case types.FrameworkTypePyTorch:
		return fmt.Sprintf("%s-%s", pod.Labels[pytorchReplicaTypeLabel], pod.Labels[pytorchReplicaIndexLabel])
	default:
		return "None"
	}
}
