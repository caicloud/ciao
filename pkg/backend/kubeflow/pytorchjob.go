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
	"github.com/caicloud/ciao/pkg/types"
)

func (b Backend) createPyTorchJob(parameter *types.Parameter) (*types.Job, error) {
	pytorchJob := b.Generator.GeneratePyTorchJob(parameter)
	pytorchJob, err := b.PyTorchJobClient.KubeflowV1alpha2().PyTorchJobs(b.Namespace).Create(pytorchJob)
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
