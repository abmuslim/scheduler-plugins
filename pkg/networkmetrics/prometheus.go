package networkmetrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

const (
	prometheusService          = "prometheus-1616380099-server"
	defaultNetworkInterface    = "ens192"
	prometheusDefaultNamespace = "monitoring"
	// highestMeasureQueryTemplate is the template string to get the query for the interface with highest used bandwidth
	highestMeasureQueryTemplate = "topk(1,sum_over_time(node_network_receive_bytes_total{device=\"%s\"}[%s]))"
	// nodeMeasureQueryTemplate is the template string to get the query for the node used bandwidth
	nodeMeasureQueryTemplate = "sum_over_time(node_network_receive_bytes_total{kubernetes_node=\"%s\",device=\"%s\"}[%s])"
)

type prometheusResponse struct {
	http.Response
}

type Prometheus struct {
	networkInterface string
	timeRange        time.Duration
	namespace        string
	api              v1.API
}

func NewDefaultPrometheus() *Prometheus {
	client, err := api.NewClient(api.Config{
		Address: fmt.Sprintf("http://%s.%s", prometheusService, prometheusDefaultNamespace),
	})
	if err != nil {
		klog.Fatalf("Error getting prometheus client: %s", err.Error())
	}

	return &Prometheus{
		networkInterface: defaultNetworkInterface,
		timeRange:        time.Minute * 5,
		namespace:        prometheusDefaultNamespace,
		api:              v1.NewAPI(client),
	}
}

func NewPrometheus(address string) *Prometheus {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		klog.Fatalf("Error getting prometheus client: %s", err.Error())
	}

	return &Prometheus{
		networkInterface: defaultNetworkInterface,
		timeRange:        time.Minute * 5,
		namespace:        prometheusDefaultNamespace,
		api:              v1.NewAPI(client),
	}
}

func (p *Prometheus) score(nodeName string) (int64, *framework.Status) {
	highestBand, err := p.getHighestBandwidthMeasure()
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("error getting highest bandwidth measure: %s", err))
	}

	nodeBandwidth, err := p.getNodeBandwidthMeasure(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("error getting node bandwidth measure: %s", err))
	}

	highestBandValue := highestBand.Value
	nodeBandwidthValue := nodeBandwidth.Value

	fmt.Printf("highestBand: %s; node bandwidth: %s", highestBandValue, nodeBandwidthValue)
	return 0, nil
}

func (p *Prometheus) getNodeBandwidthMeasure(node string) (*model.Sample, error) {
	query := fmt.Sprintf(nodeMeasureQueryTemplate, node, defaultNetworkInterface, p.timeRange)
	res, err := p.query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying prometheus: %w", err)
	}

	nodeMeasure := res.(model.Vector)
	if len(nodeMeasure) != 1 {
		return nil, fmt.Errorf("invalid response, expected 1 value, got %d", len(nodeMeasure))
	}

	return nodeMeasure[0], nil
}

func (p *Prometheus) getHighestBandwidthMeasure() (*model.Sample, error) {
	query := fmt.Sprintf(highestMeasureQueryTemplate, defaultNetworkInterface, p.timeRange)
	res, err := p.query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying prometheus: %w", err)
	}

	highestMeasure := res.(model.Vector)
	if len(highestMeasure) != 1 {
		return nil, fmt.Errorf("invalid response, expected 1 value, got %d", len(highestMeasure))
	}

	return highestMeasure[0], nil
}

func (p *Prometheus) query(query string) (model.Value, error) {
	results, warnings, err := p.api.Query(context.Background(), query, time.Now())

	if len(warnings) > 0 {
		klog.Warningf("Warnings: %v\n", warnings)
	}

	klog.Infof("result:\n%v\n", results)

	return results, err
}
