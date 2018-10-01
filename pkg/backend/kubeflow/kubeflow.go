package kubeflow

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"

	pyttorchjobclient "github.com/kubeflow/pytorch-operator/pkg/client/clientset/versioned"
	tfjobclient "github.com/kubeflow/tf-operator/pkg/client/clientset/versioned"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	restclientset "k8s.io/client-go/rest"

	"github.com/caicloud/ciao/pkg/backend/kubeflow/generator"
	"github.com/caicloud/ciao/pkg/types"
)

const (
	UserAgent     = "kubeflow-kernel"
	queryInterval = 6 * time.Second
	retryTime     = 10
)

// Backend is the type for kubeflow backend.
type Backend struct {
	TFJobClient      tfjobclient.Interface
	PyTorchJobClient pyttorchjobclient.Interface
	K8sClient        kubeclient.Interface
	Generator        generator.Interface
}

// New returns a new Backend.
func New(config *restclientset.Config) (*Backend, error) {
	tfJobClient, err := tfjobclient.NewForConfig(restclientset.AddUserAgent(config, UserAgent))
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubeclient.NewForConfig(restclientset.AddUserAgent(config, UserAgent))
	if err != nil {
		return nil, err
	}
	pytorchClient, err := pyttorchjobclient.NewForConfig(restclientset.AddUserAgent(config, UserAgent))
	if err != nil {
		return nil, err
	}

	return &Backend{
		TFJobClient:      tfJobClient,
		K8sClient:        k8sClient,
		PyTorchJobClient: pytorchClient,
		Generator:        generator.NewNative(),
	}, nil
}

func NewWithCM(config *restclientset.Config) (*Backend, error) {
	tfJobClient, err := tfjobclient.NewForConfig(restclientset.AddUserAgent(config, UserAgent))
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubeclient.NewForConfig(restclientset.AddUserAgent(config, UserAgent))
	if err != nil {
		return nil, err
	}
	pytorchClient, err := pyttorchjobclient.NewForConfig(restclientset.AddUserAgent(config, UserAgent))
	if err != nil {
		return nil, err
	}

	return &Backend{
		TFJobClient:      tfJobClient,
		K8sClient:        k8sClient,
		PyTorchJobClient: pytorchClient,
		Generator:        generator.NewCM(),
	}, nil
}

// ExecCode executes the code according to the parameter.
func (b *Backend) ExecCode(parameter *types.Parameter) (*types.Job, error) {
	switch parameter.Framework {
	case types.FrameworkTypeTensorFlow:
		return b.createTFJob(parameter)
	case types.FrameworkTypePyTorch:
		return b.createPyTorchJob(parameter)
	default:
		return nil, fmt.Errorf("Failed to get the framework %s", parameter.Framework)
	}
}

// GetLogs outputs logs for the given job.
func (b *Backend) GetLogs(job *types.Job) {
	var pods *v1.PodList
	var wg sync.WaitGroup
	var err error

	fmt.Printf("[kubeflow] Getting %s Job %s\n", job.Framework, job.Name)

	retry := 0
	for {
		pods, err = b.K8sClient.CoreV1().Pods(metav1.NamespaceDefault).List(metav1.ListOptions{
			LabelSelector: GetLabelSelectorForJob(job),
		})
		if err != nil {
			fmt.Printf("[kubeflow] Failed to get pods for the given job %s\n", job.Name)
		}
		if len(pods.Items) != job.PS+job.Worker+job.Master {
			fmt.Printf("[kubeflow] Waiting for all replicas (%d, %d, %d)\n", job.Master, job.PS, job.Worker)
			time.Sleep(queryInterval)

			if retry == retryTime {
				fmt.Printf("Tried %d times but cannot get the pods\n", retryTime)
				return
			}
			retry++
			continue
		} else {
			break
		}
	}

	fmt.Printf("[kubeflow] There are %d pods for the job.\n", len(pods.Items))

	wg.Add(len(pods.Items))
	for _, pod := range pods.Items {
		go b.getLogForPod(job, pod, wg)
	}
	wg.Wait()
	fmt.Printf("[kubeflow] Finished\n")
	return
}

func (b Backend) getLogForPod(job *types.Job, pod v1.Pod, wg sync.WaitGroup) {
	var readCloser io.ReadCloser
	var err error

	logOpts := &v1.PodLogOptions{
		Follow: true,
	}
	instanceName := GetReplicaInstanceForPod(job, pod)

	PodRef := &pod

	// Check if the pod is ready.
	for {
		if PodRef.Status.Phase == v1.PodPending {
			fmt.Printf("[kubeflow][%s] Pod is pending...\n", instanceName)
			time.Sleep(queryInterval)
			PodRef, err = b.K8sClient.CoreV1().Pods(metav1.NamespaceDefault).Get(pod.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Printf("[kubeflow][%s] Failed to get the pod\n", instanceName)
				return
			}
			continue
		}
		break
	}

	retry := 0
	for {
		readCloser, err = b.K8sClient.CoreV1().Pods(metav1.NamespaceDefault).GetLogs(PodRef.Name, logOpts).Stream()
		if err != nil {
			fmt.Printf("[kubeflow][%s] Failed to get the log of pod: %v\n", instanceName, err)
			time.Sleep(queryInterval)
			fmt.Printf("[kubeflow][%s] Retry to get log...\n", instanceName)

			if retry == retryTime {
				fmt.Printf("[kubeflow][%s] Tried %d times but cannot get the logs\n", instanceName, retryTime)
				return
			}
			retry++
			continue
		}
		break
	}
	fmt.Printf("[kubeflow][%s] Begin reading the log...\n", instanceName)

	defer readCloser.Close()
	defer wg.Done()

	reader := bufio.NewReader(readCloser)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("[kubeflow][%s] ----Terminated----\n", instanceName)
				return
			}
			fmt.Printf("[kubeflow][%s] Failed to read log: %v\n", instanceName, err)
			return
		}

		fmt.Printf("[kubeflow][%s] %s", instanceName, string(line))
	}
}
