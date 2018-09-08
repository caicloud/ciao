package simple

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	prefix       = "kubeflow-kernel-code"
	codeFile     = "code.py"
	builderImage = "caicloud/tensorflow-s2i:1.10.1-py3"
	imageOwner   = "caicloud"
)

type S2IClient struct {
}

func New() *S2IClient {
	return &S2IClient{}
}

func (s S2IClient) SourceToImage(code, jobName string) (string, error) {
	dir, err := ioutil.TempDir(os.TempDir(), prefix)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filepath.Join(dir, codeFile), []byte(code), 0666)
	if err != nil {
		return "", err
	}

	// This is a hack to let kubernetes do not pull from docker registry.
	imageName := fmt.Sprintf("%s:v1", filepath.Join(imageOwner, jobName))

	cmd := exec.Command("s2i", "build", dir, builderImage, imageName)
	cmd.Dir = dir
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	// TODO: Push to a Docker Registry.

	return imageName, err
}
