package kubeflow

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"

	tfv1alpha2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"
	kubeflowclient "github.com/kubeflow/tf-operator/pkg/client/clientset/versioned"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	restclientset "k8s.io/client-go/rest"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	namespaceDefault = "default"
	userAgent        = "kubeflow-kernel"
	queryInterval    = 6 * time.Second
	retryTime        = 10
)

type Backend struct {
	KubeflowClient kubeflowclient.Interface
	K8sClient      kubeclient.Interface
}

func New(config *restclientset.Config) (*Backend, error) {
	kubeflowClient, err := kubeflowclient.NewForConfig(restclientset.AddUserAgent(config, userAgent))
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubeclient.NewForConfig(restclientset.AddUserAgent(config, userAgent))
	if err != nil {
		return nil, err
	}

	return &Backend{
		KubeflowClient: kubeflowClient,
		K8sClient:      k8sClient,
	}, nil
}

func (b *Backend) ExecCode(parameter *types.Parameter) (*types.Job, error) {
	psCount := int32(parameter.PSCount)
	workerCount := int32(parameter.WorkerCount)

	// TODO: Using a function to generate it.
	tfJob := &tfv1alpha2.TFJob{
		TypeMeta: metav1.TypeMeta{
			Kind: tfv1alpha2.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: tfv1alpha2.TFJobSpec{
			TFReplicaSpecs: map[tfv1alpha2.TFReplicaType]*tfv1alpha2.TFReplicaSpec{
				tfv1alpha2.TFReplicaTypePS: &tfv1alpha2.TFReplicaSpec{
					Replicas: &psCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:  "tensorflow",
									Image: parameter.Image,
								},
							},
						},
					},
				},
				tfv1alpha2.TFReplicaTypeWorker: &tfv1alpha2.TFReplicaSpec{
					Replicas: &workerCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:  "tensorflow",
									Image: parameter.Image,
								},
							},
						},
					},
				},
			},
		},
	}
	tfJob, err := b.KubeflowClient.KubeflowV1alpha2().TFJobs(namespaceDefault).Create(tfJob)
	if err != nil {
		return nil, err
	}
	return &types.Job{
		Name:      tfJob.Name,
		Framework: types.FrameworkTypeTensorFlow,
		PS:        parameter.PSCount,
		Worker:    parameter.WorkerCount,
	}, nil
}

// GetLogs outputs logs for the given job.
func (b *Backend) GetLogs(job *types.Job) {
	var pods *v1.PodList
	var wg sync.WaitGroup
	var err error

	fmt.Printf("[kubeflow] Getting %s Job %s\n", job.Framework, job.Name)

	retry := 0
	for {
		pods, err = b.K8sClient.CoreV1().Pods(namespaceDefault).List(metav1.ListOptions{
			LabelSelector: GetLabelSelectorForJob(job),
		})
		if err != nil {
			fmt.Printf("[kubeflow] Failed to get pods for the given job %s\n", job.Name)
		}
		if len(pods.Items) != job.PS+job.Worker {
			fmt.Printf("[kubeflow] Wating for %d PS and %d workers\n", job.PS, job.Worker)
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
			PodRef, err = b.K8sClient.CoreV1().Pods(namespaceDefault).Get(pod.Name, metav1.GetOptions{})
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
		readCloser, err = b.K8sClient.CoreV1().Pods(namespaceDefault).GetLogs(PodRef.Name, logOpts).Stream()
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
