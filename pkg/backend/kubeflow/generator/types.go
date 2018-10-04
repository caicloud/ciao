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
