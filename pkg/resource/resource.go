package resource

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	CPU    = "cpu"
	Memory = "memory"
)

type JobResource struct {
	WorkerCPU    string
	WorkerMemory string
	PSCPU        string
	PSMemory     string
	MasterCPU    string
	MasterMemory string
}

func (res JobResource) WorkerLimits() (v1.ResourceList, error) {
	resList := v1.ResourceList{}
	if res.WorkerCPU != "" {
		cpu, err := resource.ParseQuantity(res.WorkerCPU)
		if err != nil {
			return nil, err
		}
		resList[CPU] = cpu
	}

	if res.WorkerMemory != "" {
		mem, err := resource.ParseQuantity(res.WorkerMemory)
		if err != nil {
			return nil, err
		}
		resList[Memory] = mem
	}

	return resList, nil
}

func (res JobResource) PSLimits() (v1.ResourceList, error) {
	resList := v1.ResourceList{}
	if res.PSCPU != "" {
		cpu, err := resource.ParseQuantity(res.PSCPU)
		if err != nil {
			return nil, err
		}
		resList[CPU] = cpu
	}

	if res.PSMemory != "" {
		mem, err := resource.ParseQuantity(res.PSMemory)
		if err != nil {
			return nil, err
		}
		resList[Memory] = mem
	}

	return resList, nil
}

func (res JobResource) MasterLimits() (v1.ResourceList, error) {
	resList := v1.ResourceList{}
	if res.MasterCPU != "" {
		cpu, err := resource.ParseQuantity(res.MasterCPU)
		if err != nil {
			return nil, err
		}
		resList[CPU] = cpu
	}

	if res.MasterMemory != "" {
		mem, err := resource.ParseQuantity(res.MasterMemory)
		if err != nil {
			return nil, err
		}
		resList[Memory] = mem
	}

	return resList, nil
}
