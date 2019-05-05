package resource

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"reflect"
	"testing"
)

func TestWorkerLimits(t *testing.T) {
	type TestCase struct {
		Resource JobResource
		Expected v1.ResourceList
	}
	testCases := []TestCase{
		{
			JobResource{
				WorkerCPU: "1000m",
			},
			v1.ResourceList{CPU: resource.MustParse("1000m")},
		},
		{
			JobResource{
				WorkerMemory: "1Gi",
			},
			v1.ResourceList{Memory: resource.MustParse("1Gi")},
		},
		{
			JobResource{
				WorkerCPU:    "1000m",
				WorkerMemory: "1Gi",
			},
			v1.ResourceList{
				CPU:    resource.MustParse("1000m"),
				Memory: resource.MustParse("1Gi"),
			},
		},
	}

	for _, testCase := range testCases {
		actual, err := testCase.Resource.WorkerLimits()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, testCase.Expected) {
			t.Errorf("Expected %v, got %v", testCase.Expected, actual)
		}
	}
}

func TestPSLimits(t *testing.T) {
	type TestCase struct {
		Resource JobResource
		Expected v1.ResourceList
	}
	testCases := []TestCase{
		{
			JobResource{
				PSCPU: "1000m",
			},
			v1.ResourceList{CPU: resource.MustParse("1000m")},
		},
		{
			JobResource{
				PSMemory: "1Gi",
			},
			v1.ResourceList{Memory: resource.MustParse("1Gi")},
		},
		{
			JobResource{
				PSCPU:    "1000m",
				PSMemory: "1Gi",
			},
			v1.ResourceList{
				CPU:    resource.MustParse("1000m"),
				Memory: resource.MustParse("1Gi"),
			},
		},
	}

	for _, testCase := range testCases {
		actual, err := testCase.Resource.PSLimits()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, testCase.Expected) {
			t.Errorf("Expected %v, got %v", testCase.Expected, actual)
		}
	}
}

func TestMasterLimits(t *testing.T) {
	type TestCase struct {
		Resource JobResource
		Expected v1.ResourceList
	}
	testCases := []TestCase{
		{
			JobResource{
				MasterCPU: "1000m",
			},
			v1.ResourceList{CPU: resource.MustParse("1000m")},
		},
		{
			JobResource{
				MasterMemory: "1Gi",
			},
			v1.ResourceList{Memory: resource.MustParse("1Gi")},
		},
		{
			JobResource{
				MasterCPU:    "1000m",
				MasterMemory: "1Gi",
			},
			v1.ResourceList{
				CPU:    resource.MustParse("1000m"),
				Memory: resource.MustParse("1Gi"),
			},
		},
	}

	for _, testCase := range testCases {
		actual, err := testCase.Resource.MasterLimits()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, testCase.Expected) {
			t.Errorf("Expected %v, got %v", testCase.Expected, actual)
		}
	}
}
