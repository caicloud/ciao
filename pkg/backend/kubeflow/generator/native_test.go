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

package generator

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	pytorchv1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"

	"github.com/caicloud/ciao/pkg/resource"
	"github.com/caicloud/ciao/pkg/types"
)

func TestNewTFJob(t *testing.T) {
	cm := NewNative(defaultNamespace)

	expectedPSCount := 1
	expectedWorkerCount := 1
	expectedImage := "image"
	expectedCleanPolicy := types.CleanPodPolicyAll
	expectedPSLimits := v1.ResourceList{
		resource.CPU:    k8sresource.MustParse("100m"),
		resource.Memory: k8sresource.MustParse("100Mi"),
	}
	expectedWorkerLimits := v1.ResourceList{
		resource.CPU:    k8sresource.MustParse("1000m"),
		resource.Memory: k8sresource.MustParse("1Gi"),
	}

	param := &types.Parameter{
		PSCount:     expectedPSCount,
		WorkerCount: expectedWorkerCount,
		Image:       expectedImage,
		CleanPolicy: types.CleanPodPolicyAll,
		Resource: resource.JobResource{
			WorkerCPU:    "1000m",
			WorkerMemory: "1Gi",
			PSCPU:        "100m",
			PSMemory:     "100Mi",
		},
	}

	tfJob, err := cm.GenerateTFJob(param)
	if err != nil {
		t.Fatal(err)
	}
	actualPSCount := *tfJob.Spec.TFReplicaSpecs[tfv1.TFReplicaTypePS].Replicas
	actualWorkerCount := *tfJob.Spec.TFReplicaSpecs[tfv1.TFReplicaTypeWorker].Replicas
	actualImage := tfJob.Spec.TFReplicaSpecs[tfv1.TFReplicaTypePS].Template.Spec.Containers[0].Image
	actualCleanPolicy := *tfJob.Spec.CleanPodPolicy
	actualPSLimits := tfJob.Spec.TFReplicaSpecs[tfv1.TFReplicaTypePS].
		Template.Spec.Containers[0].Resources.Limits
	actualWorkerLimits := tfJob.Spec.TFReplicaSpecs[tfv1.TFReplicaTypeWorker].
		Template.Spec.Containers[0].Resources.Limits

	if actualPSCount != int32(expectedPSCount) {
		t.Errorf("Expected %d ps, got %d", expectedPSCount, actualPSCount)
	}
	if actualWorkerCount != int32(expectedWorkerCount) {
		t.Errorf("Expected %d workers, got %d", expectedWorkerCount, actualWorkerCount)
	}
	if actualImage != expectedImage {
		t.Errorf("Expected configmap name %s, got %s", expectedImage, actualImage)
	}
	if actualCleanPolicy != common.CleanPodPolicy(expectedCleanPolicy) {
		t.Errorf("Expected clean policy %s, got %s", expectedCleanPolicy, actualCleanPolicy)
	}
	if !reflect.DeepEqual(actualPSLimits, expectedPSLimits) {
		t.Errorf("Expected ps resource limits %v, got %v", expectedPSLimits, actualPSLimits)
	}
	if !reflect.DeepEqual(actualWorkerLimits, expectedWorkerLimits) {
		t.Errorf("Expected worker resource limits %v, got %v", expectedWorkerLimits, actualWorkerLimits)
	}
}

func TestNewPyTorchJob(t *testing.T) {
	cm := NewNative(defaultNamespace)

	expectedMasterCount := 1
	expectedWorkerCount := 1
	expectedImage := "image"
	expectedCleanPolicy := types.CleanPodPolicyAll
	expectedMasterLimits := v1.ResourceList{
		resource.CPU:    k8sresource.MustParse("100m"),
		resource.Memory: k8sresource.MustParse("100Mi"),
	}
	expectedWorkerLimits := v1.ResourceList{
		resource.CPU:    k8sresource.MustParse("1000m"),
		resource.Memory: k8sresource.MustParse("1Gi"),
	}

	param := &types.Parameter{
		MasterCount: expectedMasterCount,
		WorkerCount: expectedWorkerCount,
		Image:       expectedImage,
		CleanPolicy: types.CleanPodPolicyAll,
		Resource: resource.JobResource{
			WorkerCPU:    "1000m",
			WorkerMemory: "1Gi",
			MasterCPU:    "100m",
			MasterMemory: "100Mi",
		},
	}

	pytorchJob, err := cm.GeneratePyTorchJob(param)
	if err != nil {
		t.Fatal(err)
	}
	actualMasterCount := *pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1.PyTorchReplicaTypeMaster].Replicas
	actualWorkerCount := *pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1.PyTorchReplicaTypeWorker].Replicas
	actualImage := pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1.PyTorchReplicaTypeMaster].Template.Spec.Containers[0].Image
	actualCleanPolicy := *pytorchJob.Spec.CleanPodPolicy
	actualMasterLimits := pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1.PyTorchReplicaTypeMaster].
		Template.Spec.Containers[0].Resources.Limits
	actualWorkerLimits := pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1.PyTorchReplicaTypeWorker].
		Template.Spec.Containers[0].Resources.Limits

	if actualMasterCount != int32(expectedMasterCount) {
		t.Errorf("Expected %d masters, got %d", expectedMasterCount, actualMasterCount)
	}
	if actualWorkerCount != int32(expectedWorkerCount) {
		t.Errorf("Expected %d workers, got %d", expectedWorkerCount, actualWorkerCount)
	}
	if actualImage != expectedImage {
		t.Errorf("Expected configmap name %s, got %s", expectedImage, actualImage)
	}
	if actualCleanPolicy != common.CleanPodPolicy(expectedCleanPolicy) {
		t.Errorf("Expected clean policy %s, got %s", expectedCleanPolicy, actualCleanPolicy)
	}
	if !reflect.DeepEqual(actualMasterLimits, expectedMasterLimits) {
		t.Errorf("Expected master resource limits %v, got %v", expectedMasterLimits, actualMasterLimits)
	}
	if !reflect.DeepEqual(actualWorkerLimits, expectedWorkerLimits) {
		t.Errorf("Expected worker resource limits %v, got %v", expectedWorkerLimits, actualWorkerLimits)
	}
}
