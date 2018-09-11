package types

// Parameter is the type for the parameter of the job.
type Parameter struct {
	Framework   FrameworkType
	PSCount     int
	WorkerCount int
	// MasterCount is for PyTorchJob.
	MasterCount  int
	GenerateName string
	Image        string
}

// Job is the type for general jobs (PyTorchJob, TFJob).
type Job struct {
	Framework FrameworkType
	Name      string
	PS        int
	// Master is for PyTorchJob.
	Master int
	Worker int
}

// FrameworkType is the type for types of the frameworks.
type FrameworkType string

const (
	// FrameworkTypeTensorFlow defines tensorflow type.
	FrameworkTypeTensorFlow = "tensorflow"
	// FrameworkTypePyTorch defines tensorflow type.
	FrameworkTypePyTorch = "pytorch"
)
