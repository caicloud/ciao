package mock

// Mocker is the type for S2I mocker.
type Mocker struct{}

// New returns a new Mocker.
func New() *Mocker {
	return &Mocker{}
}

// SourceToImage uses the default image without building one.
func (m Mocker) SourceToImage(code, jobName string) (string, error) {
	return "kubeflow/tf-dist-mnist-test:1.0", nil
}
