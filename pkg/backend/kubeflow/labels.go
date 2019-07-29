// Copyright 2018 Caicloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubeflow

import (
	"fmt"

	"k8s.io/api/core/v1"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	labelTFJobName      = "tf-job-name"
	tfReplicaTypeLabel  = "tf-replica-type"
	tfReplicaIndexLabel = "tf-replica-index"

	labelPyTorchJobName      = "pytorch-job-name"
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
