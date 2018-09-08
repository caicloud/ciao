package backend

import (
	"github.com/caicloud/ciao/pkg/types"
)

type Interface interface {
	ExecCode(parameter *types.Parameter) (*types.Job, error)
	GetLogs(job *types.Job)
}
