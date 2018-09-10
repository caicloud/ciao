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

// S2IClient is the type for using s2i.
type S2IClient struct {
}

// New returns a new S2IClient.
func New() *S2IClient {
	return &S2IClient{}
}

// SourceToImage converts the code to the image.
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

	// TODO(gaocegege): Push to a Docker Registry.

	return imageName, err
}
