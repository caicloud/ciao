package mock

type Mocker struct{}

func New() *Mocker {
	return &Mocker{}
}

func (m Mocker) SourceToImage(code, jobName string) (string, error) {
	return "kubeflow/tf-dist-mnist-test:1.0", nil
}
