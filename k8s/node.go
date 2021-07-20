package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/transport"
	"time"
)

const (
	HostnameLabel        = "kubernetes.io/hostname"
	NodeRoleLabel        = "node-role.kubernetes.io/master"
	MaxRetries           = 5
	RetryInterval        = 5
	WrapTransportTimeout = 30
)

// NewClient is get clientSet by kubeConfig
func NewClient(kubeConfigPath string, k8sWrapTransport transport.WrapperFunc) (*kubernetes.Clientset, error) {
	// use the current admin kubeconfig
	var config *rest.Config
	var err error
	//if home, _ := os.UserHomeDir(); home != "" && kubeConfigPath != "" {
	//	kubeConfigPath = filepath.Join(home, ".kube", "config")
	//}
	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath); err != nil {
			return nil, err
		}
	}

	if k8sWrapTransport != nil {
		config.WrapTransport = k8sWrapTransport
	}
	config.Timeout = time.Second * time.Duration(WrapTransportTimeout)
	K8sClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return K8sClientSet, nil
}
