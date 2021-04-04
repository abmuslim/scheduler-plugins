package networkmetrics

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

// NetworkTraffic is a score plugin that favors nodes based on their
// network traffic amount. Nodes with less traffic are favored.
type NetworkTraffic struct {
	handle framework.FrameworkHandle
}

// NetworkTrafficName is the name of the plugin used in the Registry and configurations.
const NetworkTrafficName = "NodeNetworkTrafficScorer"

var _ = framework.ScorePlugin(&NetworkTraffic{})

// Name returns name of the plugin. It is used in logs, etc.
func (n *NetworkTraffic) Name() string {
	return NetworkTrafficName
}

func (n *NetworkTraffic) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := n.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
	}
	fmt.Print(nodeInfo)

	return 0, nil
}

func (n *NetworkTraffic) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
