package types

// Parameter is the type for the parameter of the job.
// TODO: Add more configs.
type Parameter struct {
	Framework    FrameworkType
	PSCount      int
	WorkerCount  int
	GenerateName string
	Image        string
}

// Job is the type for general jobs (PyTorchJob, TFJob).
type Job struct {
	Framework FrameworkType
	Name      string
	PS        int
	Worker    int
}

// FrameworkType is the type for types of the frameworks.
type FrameworkType string

const (
	// FrameworkTypeTensorFlow defines tensorflow type.
	FrameworkTypeTensorFlow = "tensorflow"
)
