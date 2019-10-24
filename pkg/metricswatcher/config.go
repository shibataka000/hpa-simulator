package metricswatcher

import (
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type config struct {
	namespace                     string
	deployment                    string
	cpuInitializationPeriod       time.Duration
	delayOfInitialReadinessStatus time.Duration
	resource                      v1.ResourceName
	selector                      labels.Selector
	targetUtilization             int32
	tolerance                     float64
}

func NewConfig(namespace string, deployment string) (*config, error) {
	cpuInitializationPeriod, err := time.ParseDuration("5m")
	if err != nil {
		return nil, err
	}
	delayOfInitialReadinessStatus, err := time.ParseDuration("30s")
	if err != nil {
		return nil, err
	}

	return &config{
		namespace,
		deployment,
		cpuInitializationPeriod,
		delayOfInitialReadinessStatus,
		v1.ResourceCPU,
		labels.Everything(),
		50,
		0.1,
	}, nil
}
