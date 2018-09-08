package simple

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/caicloud/ciao/pkg/types"
)

type Interpreter struct {
	FrameworkPrefix string
	WorkerPrefix    string
	PSPrefix        string
}

func New() *Interpreter {
	return &Interpreter{
		FrameworkPrefix: "%kubeflow framework=",
		WorkerPrefix:    "%kubeflow worker=",
		PSPrefix:        "%kubeflow ps=",
	}
}

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
