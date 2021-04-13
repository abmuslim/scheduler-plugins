package networkmetrics

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"sigs.k8s.io/scheduler-plugins/pkg/apis/config"
)

// NetworkTraffic is a score plugin that favors nodes based on their
// network traffic amount. Nodes with less traffic are favored.
type NetworkTraffic struct {
	handle     framework.FrameworkHandle
	prometheus *PrometheusHandle
}

// Name is the name of the plugin used in the Registry and configurations.
const Name = "NodeNetworkTrafficScorer"

var _ = framework.ScorePlugin(&NetworkTraffic{})

// New initializes a new plugin and returns it.
func New(obj runtime.Object, h framework.FrameworkHandle) (framework.Plugin, error) {
	fmt.Printf("My custom print network traffic: %+v", obj)
	args, ok := obj.(*config.NetworkTrafficArgs)
	if !ok {
		fmt.Printf("%+v", obj)
		return nil, fmt.Errorf("my error: want args to be of type NetworkTrafficArgs, got %T", args)
	}

	fmt.Print("successfully initiated")

	return &NetworkTraffic{
		handle:     h,
		prometheus: NewPrometheus(args.Address, args.NetworkInterface, time.Minute*time.Duration(args.TimeRangeInMinutes)),
	}, nil
}

// Name returns name of the plugin. It is used in logs, etc.
func (n *NetworkTraffic) Name() string {
	return Name
}

func (n *NetworkTraffic) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeBandwidth, err := n.prometheus.getNodeBandwidthMeasure(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("error getting node bandwidth measure: %s", err))
	}

	nodeBandwidthValue := nodeBandwidth.Value

	fmt.Printf("node bandwidth: %s", nodeBandwidthValue)
	return int64(nodeBandwidth.Value), nil
}

func (n *NetworkTraffic) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

func (n *NetworkTraffic) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	//framework.MaxNodeScore
	var higherScore int64
	for _, node := range scores {
		if higherScore < node.Score {
			higherScore = node.Score
		}
	}

	for _, node := range scores {
		node.Score = framework.MaxNodeScore - (node.Score * framework.MaxNodeScore / higherScore)
	}
	return nil
}
