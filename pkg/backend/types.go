package backend

import (
	"github.com/caicloud/ciao/pkg/types"
)

// Interface is the type for backend.
type Interface interface {
	ExecCode(parameter *types.Parameter) (*types.Job, error)
	GetLogs(job *types.Job)
}
