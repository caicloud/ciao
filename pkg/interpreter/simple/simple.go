package simple

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/caicloud/ciao/pkg/types"
)

// Interpreter is the type for the simple interpreter.
type Interpreter struct {
	FrameworkPrefix string
	WorkerPrefix    string
	PSPrefix        string
}

// New returns a new interpreter.
func New() *Interpreter {
	return &Interpreter{
		FrameworkPrefix: "%framework=",
		WorkerPrefix:    "%worker=",
		PSPrefix:        "%ps=",
	}
}

// Preprocess interprets the magic commands.
func (i Interpreter) Preprocess(code string) (*types.Parameter, error) {
	param := &types.Parameter{}
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if string(line[0]) == "%" {
			if err := i.parseMagicCommand(param, line); err != nil {
				return nil, err
			}
		}
	}
	return param, nil
}

func (i Interpreter) parseMagicCommand(param *types.Parameter, line string) error {
	var err error
	if strings.Contains(line, i.FrameworkPrefix) {
		param.Framework = types.FrameworkType(line[len(i.FrameworkPrefix):])
	} else if strings.Contains(line, i.WorkerPrefix) {
		param.WorkerCount, err = strconv.Atoi(line[len(i.WorkerPrefix):])
		if err != nil {
			return err
		}
	} else if strings.Contains(line, i.PSPrefix) {
		param.PSCount, err = strconv.Atoi(line[len(i.PSPrefix):])
		if err != nil {
			return err
		}
	}
	return nil
}

// PreprocessedCode gets the preprocessed code ( the code without magic commands.)
func (i Interpreter) PreprocessedCode(code string) string {
	lines := strings.Split(code, "\n")
	res := ""
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if string(line[0]) != "%" {
			res = fmt.Sprintf("%s\n%s", res, line)
		}
	}
	return res
}
