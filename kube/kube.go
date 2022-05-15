package kube

import (
	"context"
	"io"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return rest.InClusterConfig()
	}
	kubeConfigPath := os.Getenv("KUBECONFIG")
	if kubeConfigPath == "" {
		kubeConfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	return clientcmd.BuildConfigFromFlags("", kubeConfigPath)
}

func GetClient() (*kubernetes.Clientset, error) {
	config, _ := GetConfig()
	return kubernetes.NewForConfig(config)
}

func GetService(name, namespace string) (*corev1.Service, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func GetPodLogs(name, namespace, container string) (io.ReadCloser, error) {
	config, _ := GetConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	podLogOpts := corev1.PodLogOptions{Container: container, Follow: true}
	req := clientset.CoreV1().Pods(namespace).GetLogs(name, &podLogOpts)
	return req.Stream(context.TODO())
}
