package metricswatcher

import (
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	corev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type MetricsWatcher interface {
	Start() error
}

type metricsWatcher struct {
	metricsClient *resourceclient.MetricsV1beta1Client
	podInformer   corev1.PodInformer
	config        *Config
}

func NewMetricsWatcher(clientConfig *rest.Config, config *Config) (MetricsWatcher, error) {
	metricsClient := resourceclient.NewForConfigOrDie(clientConfig)

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	podInformer := informerFactory.Core().V1().Pods()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(new interface{}) {},
		UpdateFunc: func(old interface{}, new interface{}) {},
		DeleteFunc: func(old interface{}) {},
	})
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	return &metricsWatcher{metricsClient, podInformer, config}, nil
}

func (watcher *metricsWatcher) Start() error {
	for {
		if err := watcher.watch(); err != nil {
			return err
		}
	}
}

func (watcher *metricsWatcher) watch() error {
	config := watcher.config
	selector := labels.Everything()

	metrics, _, err := getResourceMetric(watcher.metricsClient, config.resource, config.namespace, selector)
	if err != nil {
		return err
	}

	podLister := watcher.podInformer.Lister()
	podList, err := podLister.Pods(watcher.config.namespace).List(selector)
	if len(podList) == 0 {
		return fmt.Errorf("len(podList) == 0")
	}
	_, ignoredPods, _ := groupPods(podList, metrics, config.resource, config.cpuInitializationPeriod, config.delayOfInitialReadinessStatus)
	removeMetricsForPods(metrics, ignoredPods)
	requests, err := calculatePodRequests(podList, config.resource)
	if err != nil {
		return err
	}
	if len(metrics) == 0 {
		return fmt.Errorf("len(metrics) == 0")
	}
	usageRatio, _, _, err := getResourceUtilizationRatio(metrics, requests, config.targetUtilization)
	if err != nil {
		return err
	}
	log.Printf("%v\n", usageRatio)

	// metrics, err := watcher.metricsClient.PodMetricses(watcher.config.namespace).List(metav1.ListOptions{LabelSelector: selector.String()})
	// if err != nil {
	// 	return err
	// }
	//
	//
	// now := time.Now().Format("2006/01/02 15:04:05.000")
	// for _, pod := range metrics.Items {
	// 	if !watcher.config.podQuery.MatchString(pod.OAName) {
	// 		continue
	// 	}
	// 	for _, container := range pod.Containers {
	// 		fmt.Printf("%v,%v,%v,%v,%v\n", now, pod.Name, container.Name, container.Usage.Cpu(), container.Usage.Memory())
	// 	}
	// }

	log.Printf("====================\n")

	return nil
}
