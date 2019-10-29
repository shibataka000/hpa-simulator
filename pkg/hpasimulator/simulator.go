package hpasimulator

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	core "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type HpaSimulator interface {
	Start() error
}

type hpaSimulator struct {
	metricsClient *metrics.MetricsV1beta1Client
	podInformer   core.PodInformer
	config        *config
}

func NewHpaSimulator(clientConfig *rest.Config, config *config) (HpaSimulator, error) {
	metricsClient := metrics.NewForConfigOrDie(clientConfig)

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

	return &hpaSimulator{metricsClient, podInformer, config}, nil
}

func (simulator *hpaSimulator) Start() error {
	currentReplicas := int32(1)
	for {
		log.Printf("[Current Replicas] %v\n", currentReplicas)
		newReplicas, err := getResourceReplicas(simulator, simulator.config, currentReplicas)
		if err != nil {
			return err
		}
		if currentReplicas != newReplicas {
			log.Printf("[Scale] %v -> %v\n", currentReplicas, newReplicas)
			currentReplicas = newReplicas
		}
		log.Printf("========================================\n")
	}
}
