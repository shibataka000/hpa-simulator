package metricswatcher

import (
	"fmt"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
)

func Watch(config *rest.Config) {
	fmt.Println("Hello World!")

	metricsClient := resourceclient.NewForConfigOrDie(config)
	metrics, err := metricsClient.PodMetricses("").List(metav1.ListOptions{})
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("len(metrics.Items) = %v\n", len(metrics.Items))
	for _, pod := range metrics.Items {
		fmt.Printf("%v\n", pod.Name)
		for _, container := range pod.Containers {
			fmt.Printf("    %v %v\n", container.Name, container.Usage.Cpu())
		}
	}
}
