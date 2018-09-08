package interpreter

import (
	"github.com/caicloud/ciao/pkg/types"
)

type Interface interface {
	Preprocess(code string) (*types.Parameter, error)
	PreprocessedCode(code string) string
}
