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

package simple

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	prefix              = "kubeflow-kernel-code."
	codeFile            = "code.py"
	builderImageTF      = "gaocegege/tensorflow-s2i:1.10.1-py3"
	builderImagePyTorch = "gaocegege/pytorch-s2i:v0.2"
	imageOwner          = "caicloud"
)

// S2IClient is the type for using s2i.
type S2IClient struct {
}

// New returns a new S2IClient.
func New() *S2IClient {
	return &S2IClient{}
}

// SourceToImage converts the code to the image.
func (s S2IClient) SourceToImage(code string, parameter *types.Parameter) (string, error) {
	dir, err := ioutil.TempDir(os.TempDir(), prefix)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filepath.Join(dir, codeFile), []byte(code), 0666)
	if err != nil {
		return "", err
	}

	// This is a hack to let kubernetes do not pull from docker registry.
	imageName := fmt.Sprintf("%s:v1", filepath.Join(imageOwner, parameter.GenerateName))

	cmd := exec.Command("s2i", "build", dir, getBuilderImage(parameter), imageName)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("[kubeflow] Failed to build the image: %s", string(output))
		return "", err
	}

	// TODO(gaocegege): Push to a Docker Registry.

	return imageName, err
}

func getBuilderImage(parameter *types.Parameter) string {
	switch parameter.Framework {
	case types.FrameworkTypeTensorFlow:
		return builderImageTF
	case types.FrameworkTypePyTorch:
		return builderImagePyTorch
	default:
		return "-1"
	}
}

func (s S2IClient) Cleanup(parameter *types.Parameter) (error) {
	return nil
}