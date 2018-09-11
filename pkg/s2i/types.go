package s2i

import (
	"github.com/caicloud/ciao/pkg/types"
)

// Interface is the interface for s2i.
type Interface interface {
	SourceToImage(code string, parameter *types.Parameter) (string, error)
}
