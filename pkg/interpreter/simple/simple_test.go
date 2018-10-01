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
