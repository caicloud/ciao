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

package configmap

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	restclientset "k8s.io/client-go/rest"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	FileName  = "code.py"
	userAgent = "kubeflow-kernel"
)

// Client is the type for s2i client powered by configmap.
type Client struct {
	K8sClient kubeclient.Interface
	Namespace string
}

// New returns a new client.
func New(config *restclientset.Config, namespace string) (*Client, error) {
	k8sClient, err := kubeclient.NewForConfig(restclientset.AddUserAgent(config, userAgent))
	if err != nil {
		return nil, err
	}

	return &Client{
		K8sClient: k8sClient,
		Namespace: namespace,
	}, nil
}

// SourceToImage converts the code to the image.
func (c Client) SourceToImage(code string, parameter *types.Parameter) (string, error) {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: c.Namespace,
		},
		Data: map[string]string{
			FileName: code,
		},
	}

	created, err := c.K8sClient.CoreV1().ConfigMaps(c.Namespace).Create(cm)
	return created.Name, err
}
