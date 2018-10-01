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
}

// New returns a new client.
func New(config *restclientset.Config) (*Client, error) {
	k8sClient, err := kubeclient.NewForConfig(restclientset.AddUserAgent(config, userAgent))
	if err != nil {
		return nil, err
	}

	return &Client{
		K8sClient: k8sClient,
	}, nil
}

// SourceToImage converts the code to the image.
func (c Client) SourceToImage(code string, parameter *types.Parameter) (string, error) {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: metav1.NamespaceDefault,
		},
		Data: map[string]string{
			FileName: code,
		},
	}

	created, err := c.K8sClient.CoreV1().ConfigMaps(metav1.NamespaceDefault).Create(cm)
	return created.Name, err
}
