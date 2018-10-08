// Copyright 2018 Caicloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package simple

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/caicloud/ciao/pkg/types"
)

// Interpreter is the type for the simple interpreter.
type Interpreter struct {
	FrameworkPrefix string
	WorkerPrefix    string
	PSPrefix        string
	MasterPrefix    string
}

// New returns a new interpreter.
func New() *Interpreter {
	return &Interpreter{
		FrameworkPrefix: "%framework=",
		WorkerPrefix:    "kf_worker=",
		PSPrefix:        "kf_ps=",
		MasterPrefix:    "kf_master=",
	}
}

// Preprocess interprets the magic commands.
func (i Interpreter) Preprocess(code string) (*types.Parameter, error) {
	param := &types.Parameter{}
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		log.Printf("Processing line: %s\n", line)
		if string(line[0:3]) == "kf_" {
			if err := i.parseMagicCommand(param, line); err != nil {
				return nil, err
			}
		}
	}
	return param, nil
}

func (i Interpreter) parseMagicCommand(param *types.Parameter, line string) error {
	var err error
	param.Framework = types.FrameworkTypeTensorFlow
	if strings.Contains(line, i.WorkerPrefix) {
		param.WorkerCount, err = strconv.Atoi(line[len(i.WorkerPrefix):])
		if err != nil {
			return err
		}
	} else if strings.Contains(line, i.PSPrefix) {
		param.PSCount, err = strconv.Atoi(line[len(i.PSPrefix):])
		if err != nil {
			return err
		}
	} else if strings.Contains(line, i.MasterPrefix) {
		param.MasterCount, err = strconv.Atoi(line[len(i.MasterPrefix):])
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
