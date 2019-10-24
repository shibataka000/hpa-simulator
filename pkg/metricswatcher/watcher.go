package metricswatcher

import (
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
	config        *config
}

func NewMetricsWatcher(clientConfig *rest.Config, config *config) (MetricsWatcher, error) {
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
	currentReplicas := int32(1)
	selector := labels.Everything()
	for {
		newReplicas, err := getResourceReplicas(watcher, watcher.config, selector, currentReplicas)
		if err != nil {
			return err
		}
		if currentReplicas != newReplicas {
			log.Printf("[Scale] %v -> %v\n", currentReplicas, newReplicas)
		}
		log.Printf("====================\n")
	}
}
