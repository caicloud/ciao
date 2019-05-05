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
	"reflect"
	"testing"

	"github.com/caicloud/ciao/pkg/resource"
	"github.com/caicloud/ciao/pkg/types"
)

// TestPreprocess test the preprocess logic.
func TestPreprocess(t *testing.T) {
	defaultRes := resource.JobResource{
		WorkerCPU:    "1000m",
		WorkerMemory: "1Gi",
		PSCPU:        "1000m",
		PSMemory:     "1Gi",
		MasterCPU:    "1000m",
		MasterMemory: "1Gi",
	}

	i := New(defaultRes)
	type TestCase struct {
		Code     string
		Expected *types.Parameter
	}
	testCases := []TestCase{
		{
			Code: `%framework=tensorflow
%ps=1
%worker=1
some code here.
`,
			Expected: &types.Parameter{
				Framework:   types.FrameworkTypeTensorFlow,
				PSCount:     1,
				WorkerCount: 1,
				Resource:    defaultRes,
			},
		},
		{
			Code: `%framework=tensorflow
%ps=1
some code here.
`,
			Expected: &types.Parameter{
				Framework: types.FrameworkTypeTensorFlow,
				PSCount:   1,
				Resource:  defaultRes,
			},
		},
		{
			Code: `%framework=tensorflow
%ps=1`,
			Expected: &types.Parameter{
				Framework: types.FrameworkTypeTensorFlow,
				PSCount:   1,
				Resource:  defaultRes,
			},
		},
		{
			Code: `%framework=tensorflow
%cleanPolicy=running`,
			Expected: &types.Parameter{
				Framework:   types.FrameworkTypeTensorFlow,
				CleanPolicy: types.CleanPodPolicyRunning,
				Resource:    defaultRes,
			},
		},
		{
			Code: `%framework=tensorflow
%cleanPolicy=all`,
			Expected: &types.Parameter{
				Framework:   types.FrameworkTypeTensorFlow,
				CleanPolicy: types.CleanPodPolicyAll,
				Resource:    defaultRes,
			},
		},
		{
			Code: `%framework=tensorflow
%cleanPolicy=none`,
			Expected: &types.Parameter{
				Framework:   types.FrameworkTypeTensorFlow,
				CleanPolicy: types.CleanPodPolicyNone,
				Resource:    defaultRes,
			},
		},
		{
			// Invalid clean policy will use default value none.
			Code: `%framework=tensorflow
%cleanPolicy=test`,
			Expected: &types.Parameter{
				Framework:   types.FrameworkTypeTensorFlow,
				CleanPolicy: types.CleanPodPolicyNone,
				Resource:    defaultRes,
			},
		},
		{
			Code: `%framework=tensorflow
%ps=1;%cpu=100m;%memory=100Mi`,
			Expected: &types.Parameter{
				Framework: types.FrameworkTypeTensorFlow,
				PSCount:   1,
				Resource: resource.JobResource{
					PSCPU:        "100m",
					PSMemory:     "100Mi",
					WorkerCPU:    defaultRes.WorkerCPU,
					WorkerMemory: defaultRes.WorkerMemory,
					MasterCPU:    defaultRes.MasterCPU,
					MasterMemory: defaultRes.MasterMemory,
				},
			},
		},
		{
			Code: `%framework=tensorflow
%ps=1;%cpu=100m;%memory=100Mi
%worker=2;%cpu=10m;%memory=10Mi`,
			Expected: &types.Parameter{
				Framework:   types.FrameworkTypeTensorFlow,
				PSCount:     1,
				WorkerCount: 2,
				Resource: resource.JobResource{
					PSCPU:        "100m",
					PSMemory:     "100Mi",
					WorkerCPU:    "10m",
					WorkerMemory: "10Mi",
					MasterCPU:    defaultRes.MasterCPU,
					MasterMemory: defaultRes.MasterMemory,
				},
			},
		},
		{
			Code: `%framework=pytorch
%master=1;%cpu=100m;%memory=100Mi
%worker=2;%cpu=10m;%memory=10Mi`,
			Expected: &types.Parameter{
				Framework:   types.FrameworkTypePyTorch,
				MasterCount: 1,
				WorkerCount: 2,
				Resource: resource.JobResource{
					PSCPU:        defaultRes.PSCPU,
					PSMemory:     defaultRes.PSMemory,
					WorkerCPU:    "10m",
					WorkerMemory: "10Mi",
					MasterCPU:    "100m",
					MasterMemory: "100Mi",
				},
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
	i := New(resource.JobResource{})
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
