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
	"testing"

	pytorchv1beta2 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1beta2"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1beta2"
	tfv1beta2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1beta2"

	"github.com/caicloud/ciao/pkg/types"
)

func TestNewTFJob(t *testing.T) {
	cm := NewNative(defaultNamespace)

	expectedPSCount := 1
	expectedWorkerCount := 1
	expectedImage := "image"
	expectedCleanPolicy := types.CleanPodPolicyAll

	param := &types.Parameter{
		PSCount:     expectedPSCount,
		WorkerCount: expectedWorkerCount,
		Image:       expectedImage,
		CleanPolicy: types.CleanPodPolicyAll,
	}

	tfJob := cm.GenerateTFJob(param)
	actualPSCount := *tfJob.Spec.TFReplicaSpecs[tfv1beta2.TFReplicaTypePS].Replicas
	actualWorkerCount := *tfJob.Spec.TFReplicaSpecs[tfv1beta2.TFReplicaTypeWorker].Replicas
	actualImage := tfJob.Spec.TFReplicaSpecs[tfv1beta2.TFReplicaTypePS].Template.Spec.Containers[0].Image
	actualCleanPolicy := *tfJob.Spec.CleanPodPolicy
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
}

func TestNewPyTorchJob(t *testing.T) {
	cm := NewNative(defaultNamespace)

	expectedMasterCount := 1
	expectedWorkerCount := 1
	expectedImage := "image"
	expectedCleanPolicy := types.CleanPodPolicyAll

	param := &types.Parameter{
		MasterCount: expectedMasterCount,
		WorkerCount: expectedWorkerCount,
		Image:       expectedImage,
		CleanPolicy: types.CleanPodPolicyAll,
	}

	pytorchJob := cm.GeneratePyTorchJob(param)
	actualMasterCount := *pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1beta2.PyTorchReplicaTypeMaster].Replicas
	actualWorkerCount := *pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1beta2.PyTorchReplicaTypeWorker].Replicas
	actualImage := pytorchJob.Spec.PyTorchReplicaSpecs[pytorchv1beta2.PyTorchReplicaTypeMaster].Template.Spec.Containers[0].Image
	actualCleanPolicy := *pytorchJob.Spec.CleanPodPolicy
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
}
