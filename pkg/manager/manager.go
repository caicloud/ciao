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

package manager

import (
	"fmt"

	"github.com/caicloud/ciao/pkg/backend"
	"github.com/caicloud/ciao/pkg/interpreter"
	"github.com/caicloud/ciao/pkg/s2i"
	"github.com/caicloud/ciao/pkg/types"
)

const (
	jobNamePrefix = "jupyter-kernel"
)

// Manager is the type for the manager.
type Manager struct {
	Backend     backend.Interface
	S2IClient   s2i.Interface
	Interpreter interpreter.Interface
}

// New creates a new Manager.
func New(Backend backend.Interface, S2IClient s2i.Interface, Interpreter interpreter.Interface) *Manager {
	return &Manager{
		Backend:     Backend,
		S2IClient:   S2IClient,
		Interpreter: Interpreter,
	}
}

// Execute executes the code.
func (m Manager) Execute(code string) (*types.Job, error) {
	parameter, err := m.Interpreter.Preprocess(code)
	if err != nil {
		return nil, err
	}

	// Generate random name for the TFJob.
	parameter.GenerateName = fmt.Sprintf("%s-%s", jobNamePrefix, RandStringRunes(5))

	preprocessedCode := m.Interpreter.PreprocessedCode(code)
	// Build and get the image from source code.
	image, err := m.GetImage(preprocessedCode, parameter)
	if err != nil {
		return nil, err
	}

	parameter.Image = image

	job, err := m.Backend.ExecCode(parameter)
	if err != nil {
		return nil, err
	}
	m.Backend.GetLogs(job)
	return job, nil
}

// GetImage gets the image by the given code.
func (m Manager) GetImage(code string, parameter *types.Parameter) (string, error) {
	fmt.Println("[kubeflow] Building the Docker image...")
	imageName, err := m.S2IClient.SourceToImage(code, parameter)
	if err != nil {
		return "", err
	}
	fmt.Println("[kubeflow] Image built successfully")
	return imageName, nil
}
