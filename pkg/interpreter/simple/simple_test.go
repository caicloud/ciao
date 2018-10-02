// Copyright [yyyy] [name of copyright owner]
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
	"reflect"
	"testing"

	"github.com/caicloud/ciao/pkg/types"
)

// TestPreprocess test the preprocess logic.
func TestPreprocess(t *testing.T) {
	i := New()
	type TestCase struct {
		Code     string
		Expected *types.Parameter
	}
	testCases := []TestCase{
		TestCase{
			Code: `%framework=tensorflow
%ps=1
%worker=1
some code here.
`,
			Expected: &types.Parameter{
				Framework:   types.FrameworkTypeTensorFlow,
				PSCount:     1,
				WorkerCount: 1,
			},
		},
		TestCase{
			Code: `%framework=tensorflow
%ps=1
some code here.
`,
			Expected: &types.Parameter{
				Framework: types.FrameworkTypeTensorFlow,
				PSCount:   1,
			},
		},
		TestCase{
			Code: `%framework=tensorflow
%ps=1`,
			Expected: &types.Parameter{
				Framework: types.FrameworkTypeTensorFlow,
				PSCount:   1,
			},
		},
	}

	for _, tc := range testCases {
		actual, err := i.Preprocess(tc.Code)
		if err != nil {
			t.Errorf("Expected nil got error: %v", err)
		}
		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Errorf("Expected %v, got %v", tc.Expected, actual)
		}
	}
}

// TestGetPreProcessedCode tests the logic about getting the preprocessed code.
func TestGetPreProcessedCode(t *testing.T) {
	i := New()
	type TestCase struct {
		Code     string
		Expected string
	}
	testCases := []TestCase{
		TestCase{
			Code: `%framework=tensorflow
%ps=1
%worker=1
some code here.
`,
			Expected: `
some code here.`,
		},
	}
	for _, tc := range testCases {
		actual := i.PreprocessedCode(tc.Code)
		if actual != tc.Expected {
			t.Errorf("Expected %s, got %s", tc.Expected, actual)
		}
	}
}
