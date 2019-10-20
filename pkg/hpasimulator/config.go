package hpasimulator

import (
	"time"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// Config object have some simulation parameters.
type config struct {
	namespace                     string
	cpuInitializationPeriod       time.Duration
	delayOfInitialReadinessStatus time.Duration
	resource                      core.ResourceName
	selector                      labels.Selector
	targetUtilization             int32
	tolerance                     float64
}

// NewConfig return configuration object.
func NewConfig(namespace string, selectorString string) (*config, error) {
	cpuInitializationPeriod, err := time.ParseDuration("5m")
	if err != nil {
		return nil, err
	}
	delayOfInitialReadinessStatus, err := time.ParseDuration("30s")
	if err != nil {
		return nil, err
	}

	selector, err := labels.Parse(selectorString)
	if err != nil {
		return nil, err
	}

	return &config{
		namespace,
		cpuInitializationPeriod,
		delayOfInitialReadinessStatus,
		core.ResourceCPU,
		selector,
		50,
		0.1,
	}, nil
}
