package networkmetrics

import "time"

const (
	prometheusPort             = ""
	prometheusService          = ""
	prometheusDefaultNamespace = ""
)

type Prometheus struct {
	networkInterface string
	timeRange        time.Duration
	namespace        string
}
