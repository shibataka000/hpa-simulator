package metricswatcher

import (
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type Config struct {
	namespace                     string
	deployment                    string
	cpuInitializationPeriod       time.Duration
	delayOfInitialReadinessStatus time.Duration
	resource                      v1.ResourceName
	selector                      labels.Selector
	targetUtilization             int32
}

func NewConfig(namespace string, deployment string) (*Config, error) {
	cpuInitializationPeriod, err := time.ParseDuration("5m")
	if err != nil {
		return nil, err
	}
	delayOfInitialReadinessStatus, err := time.ParseDuration("30s")
	if err != nil {
		return nil, err
	}

	return &Config{
		namespace,
		deployment,
		cpuInitializationPeriod,
		delayOfInitialReadinessStatus,
		v1.ResourceCPU,
		labels.Everything(),
		50,
	}, nil
}
