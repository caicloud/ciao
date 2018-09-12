package mock

import "github.com/caicloud/ciao/pkg/types"

// Mocker is the type for S2I mocker.
type Mocker struct{}

// New returns a new Mocker.
func New() *Mocker {
	return &Mocker{}
}

// SourceToImage uses the default image without building one.
func (m Mocker) SourceToImage(code string, parameter *types.Parameter) (string, error) {
	return "kubeflow/tf-dist-mnist-test:1.0", nil
}
