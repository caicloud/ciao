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

func (b Backend) createTFJob(parameter *types.Parameter) (*types.Job, error) {
	tfJob := b.Generator.GenerateTFJob(parameter)
	tfJob, err := b.TFJobClient.KubeflowV1alpha2().TFJobs(b.Namespace).Create(tfJob)
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
