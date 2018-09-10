package interpreter

import (
	"github.com/caicloud/ciao/pkg/types"
)

// Interface is the interface of the interpreter.
type Interface interface {
	Preprocess(code string) (*types.Parameter, error)
	PreprocessedCode(code string) string
}
