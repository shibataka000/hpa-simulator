package kubernetes

import (
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClientConfig(kubeconfigPath string, contextName string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{CurrentContext: contextName}).ClientConfig()
}
