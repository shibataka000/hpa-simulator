package metricswatcher

import (
	"fmt"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

// https://github.com/kubernetes/kubernetes/blob/81a8b9804a3bf310cb74446bd8a8e08e0317de22/pkg/controller/podautoscaler/metrics/interfaces.go#L28
type PodMetric struct {
	Timestamp time.Time
	Window    time.Duration
	Value     int64
}
type PodMetricsInfo map[string]PodMetric

// https://github.com/kubernetes/kubernetes/blob/81a8b9804a3bf310cb74446bd8a8e08e0317de22/pkg/controller/podautoscaler/replica_calculator.go#L62
func getResourceMetric(metricsClient *resourceclient.MetricsV1beta1Client, resource v1.ResourceName, namespace string, selector labels.Selector) (PodMetricsInfo, time.Time, error) {
	metrics, err := metricsClient.PodMetricses(namespace).List(metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("unable to fetch metrics from resource metrics API: %v", err)
	}

	if len(metrics.Items) == 0 {
		return nil, time.Time{}, fmt.Errorf("no metrics returned from resource metrics API")
	}

	res := make(PodMetricsInfo, len(metrics.Items))

	for _, m := range metrics.Items {
		podSum := int64(0)
		missing := len(m.Containers) == 0
		for _, c := range m.Containers {
			resValue, found := c.Usage[v1.ResourceName(resource)]
			if !found {
				missing = true
				log.Printf("missing resource metric %v for container %s in pod %s/%s\n", resource, c.Name, namespace, m.Name)
				break // containers loop
			}
			log.Printf("[Container Metrics] %v %v: %v\n", m.Name, c.Name, resValue.MilliValue())
			podSum += resValue.MilliValue()
		}

		if !missing {
			res[m.Name] = PodMetric{
				Timestamp: m.Timestamp.Time,
				Window:    m.Window.Duration,
				Value:     int64(podSum),
			}
		}
	}

	timestamp := metrics.Items[0].Timestamp.Time

	return res, timestamp, nil
}

// https://github.com/kubernetes/kubernetes/blob/81a8b9804a3bf310cb74446bd8a8e08e0317de22/pkg/controller/podautoscaler/replica_calculator.go#L361
func groupPods(pods []*v1.Pod, metrics PodMetricsInfo, resource v1.ResourceName, cpuInitializationPeriod, delayOfInitialReadinessStatus time.Duration) (readyPodCount int, ignoredPods sets.String, missingPods sets.String) {
	missingPods = sets.NewString()
	ignoredPods = sets.NewString()
	for _, pod := range pods {
		if pod.DeletionTimestamp != nil || pod.Status.Phase == v1.PodFailed {
			continue
		}
		// Pending pods are ignored.
		if pod.Status.Phase == v1.PodPending {
			ignoredPods.Insert(pod.Name)
			continue
		}
		// Pods missing metrics.
		metric, found := metrics[pod.Name]
		if !found {
			missingPods.Insert(pod.Name)
			continue
		}
		// Unready pods are ignored.
		if resource == v1.ResourceCPU {
			var ignorePod bool
			_, condition := podutil.GetPodCondition(&pod.Status, v1.PodReady)
			if condition == nil || pod.Status.StartTime == nil {
				ignorePod = true
			} else {
				// Pod still within possible initialisation period.
				if pod.Status.StartTime.Add(cpuInitializationPeriod).After(time.Now()) {
					// Ignore sample if pod is unready or one window of metric wasn't collected since last state transition.
					ignorePod = condition.Status == v1.ConditionFalse || metric.Timestamp.Before(condition.LastTransitionTime.Time.Add(metric.Window))
				} else {
					// Ignore metric if pod is unready and it has never been ready.
					ignorePod = condition.Status == v1.ConditionFalse && pod.Status.StartTime.Add(delayOfInitialReadinessStatus).After(condition.LastTransitionTime.Time)
				}
			}
			if ignorePod {
				ignoredPods.Insert(pod.Name)
				continue
			}
		}
		readyPodCount++
	}
	return
}

// https://github.com/kubernetes/kubernetes/blob/81a8b9804a3bf310cb74446bd8a8e08e0317de22/pkg/controller/podautoscaler/replica_calculator.go#L421
func removeMetricsForPods(metrics PodMetricsInfo, pods sets.String) {
	for _, pod := range pods.UnsortedList() {
		delete(metrics, pod)
	}
}

// https://github.com/kubernetes/kubernetes/blob/81a8b9804a3bf310cb74446bd8a8e08e0317de22/pkg/controller/podautoscaler/replica_calculator.go#L405
func calculatePodRequests(pods []*v1.Pod, resource v1.ResourceName) (map[string]int64, error) {
	requests := make(map[string]int64, len(pods))
	for _, pod := range pods {
		podSum := int64(0)
		for _, container := range pod.Spec.Containers {
			if containerRequest, ok := container.Resources.Requests[resource]; ok {
				podSum += containerRequest.MilliValue()
			} else {
				return nil, fmt.Errorf("missing request for %s", resource)
			}
		}
		requests[pod.Name] = podSum
	}
	return requests, nil
}

// https://github.com/kubernetes/kubernetes/blob/81a8b9804a3bf310cb74446bd8a8e08e0317de22/pkg/controller/podautoscaler/metrics/utilization.go#L26
func getResourceUtilizationRatio(metrics PodMetricsInfo, requests map[string]int64, targetUtilization int32) (utilizationRatio float64, currentUtilization int32, rawAverageValue int64, err error) {
	metricsTotal := int64(0)
	requestsTotal := int64(0)
	numEntries := 0

	for podName, metric := range metrics {
		request, hasRequest := requests[podName]
		if !hasRequest {
			// we check for missing requests elsewhere, so assuming missing requests == extraneous metrics
			continue
		}

		metricsTotal += metric.Value
		requestsTotal += request
		numEntries++

		log.Printf("[Pod Resource Utilization] %v: %v / %v = %v\n", podName, metric.Value, request, float64(metric.Value)/float64(request))
	}

	// if the set of requests is completely disjoint from the set of metrics,
	// then we could have an issue where the requests total is zero
	if requestsTotal == 0 {
		return 0, 0, 0, fmt.Errorf("no metrics returned matched known pods")
	}

	currentUtilization = int32((metricsTotal * 100) / requestsTotal)
	log.Printf("[Deployment Resource Utilization] %v\n", currentUtilization)

	return float64(currentUtilization) / float64(targetUtilization), currentUtilization, metricsTotal / int64(numEntries), nil
}
