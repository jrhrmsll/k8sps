package k8s

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/jrhrmsll/k8sps/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient struct {
	clientset *kubernetes.Clientset
	nodeName  string
}

func NewKubernetesClient() (*KubernetesClient, error) {
	config, err := rest.InClusterConfig()

	if err != nil {
		// fallback to kubeconfig
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		kubeconfig := path.Join(home, ".kube/config")
		if env, ok := os.LookupEnv("KUBECONFIG"); ok {
			kubeconfig = env
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("the kubeconfig cannot be loaded: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	kubernetesClient := &KubernetesClient{
		clientset: clientset,
	}

	return kubernetesClient, nil
}

func (ks *KubernetesClient) NodesIPAddresses() (map[string][]string, error) {
	nodes, err := ks.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaces, err := ks.clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	entries := make(map[string][]string)
	for _, node := range nodes.Items {
		ips := make(map[string]struct{})

		listOptions := metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.nodeName==%s", node.Name),
		}

		for _, namespace := range namespaces.Items {
			pods, err := ks.clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), listOptions)
			if err != nil {
				return nil, err
			}

			for _, pod := range pods.Items {
				ips[pod.Status.PodIP] = struct{}{}
				ips[pod.Status.HostIP] = struct{}{}
			}
		}

		entries[node.Name] = util.MapToSlice(ips)
	}

	return entries, nil
}
