package main

import (
	metricswatcher "github.com/shibataka000/metrics-watcher/pkg/metricswatcher"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"flag"
	"os"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	metricswatcher.Watch(config);
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
