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

package img

const (
	tensorflowTemplate = `FROM tensorflow/tensorflow:1.10.1-py3
LABEL maintainer="Kubeflow"
COPY ./code.py /
ENTRYPOINT ["python", "/code.py"]`

	pytorchTemplate = `FROM pytorch/pytorch:v0.2
LABEL maintainer="Kubeflow"
COPY ./code.py /
ENTRYPOINT ["python", "/code.py"]`
)
