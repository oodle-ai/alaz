package datastore

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ddosify/alaz/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/puzpuzpuz/xsync/v3"
)

var isS3PathRegex *regexp.Regexp

func init() {
	// Matches s3.amazonaws.com or bucket.s3.us-west-2.amazonaws.com
	isS3PathRegex = regexp.MustCompile(`^.*s3.*amazonaws\.com$`)
}

type OodleDS struct {
	podUidToContainerMap *xsync.MapOf[string, []Container]
	serviceUidToNameMap  *xsync.MapOf[string, string]
	metrics              *oodleMetrics
}

func NewOodleDS(reg *prometheus.Registry) *OodleDS {
	return &OodleDS{
		podUidToContainerMap: xsync.NewMapOf[string, []Container](),
		serviceUidToNameMap:  xsync.NewMapOf[string, string](),
		metrics:              newOodleMetrics(reg),
	}
}

type oodleMetrics struct {
	requestDurationMs *prometheus.HistogramVec
}

func newOodleMetrics(reg *prometheus.Registry) *oodleMetrics {
	requestDurationMs := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "oodle_ebpf_request_duration_ms",
		Help:    "The duration of requests in milliseconds",
		Buckets: []float64{100, 200, 500, 1000, 2000, 5000, 10000},
	}, []string{"path", "method", "protocol", "source", "target", "target_type", "status", "fail_reason"})

	reg.MustRegister(requestDurationMs)
	return &oodleMetrics{
		requestDurationMs: requestDurationMs,
	}
}

func (o *OodleDS) PersistPod(pod Pod, eventType string) error {
	switch eventType {
	case "DELETE":
		o.podUidToContainerMap.Delete(pod.UID)
	}
	return nil
}

func (o *OodleDS) PersistService(service Service, eventType string) error {
	switch eventType {
	case "ADD":
		o.serviceUidToNameMap.Store(service.UID, service.Name)
	case "UPDATE":
		o.serviceUidToNameMap.Store(service.UID, service.Name)
	case "DELETE":
		o.serviceUidToNameMap.Delete(service.UID)
	}
	return nil
}

func (o *OodleDS) PersistReplicaSet(rs ReplicaSet, eventType string) error {
	return nil
}

func (o *OodleDS) PersistDeployment(d Deployment, eventType string) error {
	return nil
}

func (o *OodleDS) PersistEndpoints(e Endpoints, eventType string) error {
	return nil
}

func (o *OodleDS) PersistContainer(c Container, eventType string) error {
	o.podUidToContainerMap.Compute(
		c.PodUID,
		func(oldValue []Container, loaded bool) ([]Container, bool) {
			if loaded {
				return append(oldValue, c), false
			}
			return []Container{c}, false
		},
	)
	return nil
}

func (o *OodleDS) PersistDaemonSet(ds DaemonSet, eventType string) error {
	return nil
}

func (o *OodleDS) PersistStatefulSet(ss StatefulSet, eventType string) error {
	return nil
}

func (o *OodleDS) PersistRequest(request *Request) error {
	sourceName, sOk := o.getSourceName(request)
	if !sOk {
		sourceName = fmt.Sprintf("%s:%d", request.FromIP, request.FromPort)
	}
	targetName, tOk := o.getTargetName(request)
	if !tOk {
		targetName = fmt.Sprintf("%s:%d", request.ToIP, request.ToPort)
	}

	path := request.Path

	if request.Protocol != "HTTP" && request.Protocol != "HTTPS" && request.Protocol != "gRPC" {
		// Only emit metrics for HTTP/gRPC requests. We receive DB requests as well
		// that need to be normalized if we emit them as metrics.
		return nil
	}

	// Remove query params from the path if it is a URL
	// E.g. /api/v1/namespaces/default/pods?labelSelector=app%3Dnginx => /api/v1/namespaces/default/pods
	if strings.Contains(path, "?") && strings.HasPrefix(path, "/") {
		path = strings.SplitN(path, "?", 2)[0]
	}

	// request.ToUID is set to DNS host name if the request is to a 3rd party URL
	// outside the cluster.
	if isS3PathRegex.MatchString(request.ToUID) {
		// For S3 requests, strip the bucket name from the path to avoid high cardinality metrics
		// as path is the object name for S3 calls.
		// Can use a fingerprint mechanism here to make it generic.
		// S3 Path is of the form: /<bucket-name>/<path>
		parts := strings.SplitN(path, "/", 3)
		if len(parts) >= 3 {
			path = "/" + parts[1] // Keep bucket name
		}
	}

	defer func() {
		// Given invalid labels are going to be less common case,
		// we avoid checking overhead in each request, but rather pay
		// the expensive penalty of recovering from panic in invalid requests.
		if r := recover(); r != nil {
			log.Logger.Warn().Msgf("Recovered from panic while recording metrics for request: %+v, error: %v", request, r)
		}
	}()
	o.metrics.requestDurationMs.WithLabelValues(
		path,
		request.Method,
		request.Protocol,
		sourceName,
		targetName,
		request.ToType,
		strconv.Itoa(int(request.StatusCode)),
		request.FailReason,
	).Observe(float64(request.Latency / 1e6))

	return nil
}

func (o *OodleDS) PersistKafkaEvent(request *KafkaEvent) error {
	return nil
}

func (o *OodleDS) PersistAliveConnection(trace *AliveConnection) error {
	return nil
}

func (o *OodleDS) getSourceName(request *Request) (string, bool) {
	containers, exists := o.podUidToContainerMap.Load(request.FromUID)
	if !exists {
		return "", false
	}

	// This is not accurate always. If there are multiple containers
	// in a pod, then we will always add first container in the source name.
	// To be fixed.
	return containers[0].Name, true
}

func (o *OodleDS) getTargetName(request *Request) (string, bool) {
	switch request.ToType {
	case "pod":
		containers, exists := o.podUidToContainerMap.Load(request.ToUID)
		if !exists {
			return "", false
		}

		for _, container := range containers {
			for _, p := range container.Ports {
				if p.Port == int32(request.ToPort) {
					return container.Name, true
				}
			}
		}

		return containers[0].Name, true

	case "service":
		name, exists := o.serviceUidToNameMap.Load(request.ToUID)
		if !exists {
			return "", false
		}
		return name, true

	case "outbound":
		return request.ToUID, true
	}

	return "", false
}
