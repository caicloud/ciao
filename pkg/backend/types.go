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

package backend

import (
	"github.com/caicloud/ciao/pkg/types"
)

// Interface is the type for backend.
type Interface interface {
	// ExecCode runs the job according to the parameters.
	ExecCode(parameter *types.Parameter) (*types.Job, error)
	// GetLogs get the logs from the job and output to STDOUT.
	// Outputs in STDOUT will be redirected to Jupyter Notebook's output handler.
	GetLogs(job *types.Job)
}
