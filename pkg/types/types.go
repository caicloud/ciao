package types

// TODO: Add more configs.
type Parameter struct {
	Framework    FrameworkType
	PSCount      int
	WorkerCount  int
	GenerateName string
	Image        string
}

type Job struct {
	Framework FrameworkType
	Name      string
	PS        int
	Worker    int
}

type FrameworkType string

const (
	FrameworkTypeTensorFlow = "tensorflow"
)
