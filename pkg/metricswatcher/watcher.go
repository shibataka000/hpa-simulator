package metricswatcher

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

func Watch(config *rest.Config, namespace string) error {
	metricsClient := resourceclient.NewForConfigOrDie(config)
	for {
		calcReplicas(metricsClient, namespace)
	}
}

func calcReplicas(metricsClient *resourceclient.MetricsV1beta1Client, namespace string) error {
	metrics, err := metricsClient.PodMetricses(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	now := time.Now().Format("2006/01/02 15:04:05.000")
	for _, pod := range metrics.Items {
		for _, container := range pod.Containers {
			fmt.Printf("%v,%v,%v,%v,%v\n", now, pod.Name, container.Name, container.Usage.Cpu(), container.Usage.Memory())
		}
	}
	return nil
}
