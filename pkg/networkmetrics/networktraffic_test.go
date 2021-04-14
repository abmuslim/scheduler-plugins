package networkmetrics

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/scheduler-plugins/pkg/apis/config"
	pluginconfig "sigs.k8s.io/scheduler-plugins/pkg/apis/config"
)

// type testSharedLister struct {
// 	nodes       []*v1.Node
// 	nodeInfos   []*framework.NodeInfo
// 	nodeInfoMap map[string]*framework.NodeInfo
// }

// func (f *testSharedLister) NodeInfos() framework.NodeInfoLister {
// 	return f
// }

// func (f *testSharedLister) List() ([]*framework.NodeInfo, error) {
// 	return f.nodeInfos, nil
// }

// func (f *testSharedLister) HavePodsWithAffinityList() ([]*framework.NodeInfo, error) {
// 	return nil, nil
// }

// func (f *testSharedLister) HavePodsWithRequiredAntiAffinityList() ([]*framework.NodeInfo, error) {
// 	return nil, nil
// }

// func (f *testSharedLister) Get(nodeName string) (*framework.NodeInfo, error) {
// 	return f.nodeInfoMap[nodeName], nil
// }

func TestNew(t *testing.T) {
	networkTrafficArgs := &pluginconfig.NetworkTrafficArgs{
		Address:            "localhost",
		NetworkInterface:   "ens192",
		TimeRangeInMinutes: 5,
	}

	obj := runtime.Object(networkTrafficArgs)

	args, ok := obj.(*config.NetworkTrafficArgs)
	if !ok {
		t.Fail()
	}

	fmt.Print(args)
}

// 	networkTrafficConfig := config.PluginConfig{
// 		Name: Name,
// 		Args: &networkTrafficArgs,
// 	}

// 	registeredPlugins := []st.RegisterPluginFunc{
// 		st.RegisterScorePlugin(Name, New, 1),
// 	}

// 	cs := testclientset.NewSimpleClientset()
// 	informerFactory := informers.NewSharedInformerFactory(cs, 0)
// 	snapshot := newTestSharedLister(nil, nil)

// 	fh, err := NewFramework(registeredPlugins, []config.PluginConfig{networkTrafficConfig}, runtime.WithClientSet(cs),
// 		runtime.WithInformerFactory(informerFactory), runtime.WithSnapshotSharedLister(snapshot))

// 	assert.Nil(t, err)
// 	p, err := New(&networkTrafficArgs, fh)
// 	assert.NotNil(t, p)
// 	assert.Nil(t, err)
// }

// func NewFramework(fns []st.RegisterPluginFunc, args []config.PluginConfig, opts ...runtime.Option) (framework.Framework, error) {
// 	registry := runtime.Registry{}
// 	plugins := &config.Plugins{}
// 	var pluginConfigs []config.PluginConfig
// 	for _, f := range fns {
// 		f(&registry, plugins, pluginConfigs)
// 	}
// 	return runtime.NewFramework(registry, plugins, args, opts...)
// }

// func newTestSharedLister(pods []*v1.Pod, nodes []*v1.Node) *testSharedLister {
// 	nodeInfoMap := make(map[string]*framework.NodeInfo)
// 	nodeInfos := make([]*framework.NodeInfo, 0)
// 	for _, pod := range pods {
// 		nodeName := pod.Spec.NodeName
// 		if _, ok := nodeInfoMap[nodeName]; !ok {
// 			nodeInfoMap[nodeName] = framework.NewNodeInfo()
// 		}
// 		nodeInfoMap[nodeName].AddPod(pod)
// 	}
// 	for _, node := range nodes {
// 		if _, ok := nodeInfoMap[node.Name]; !ok {
// 			nodeInfoMap[node.Name] = framework.NewNodeInfo()
// 		}
// 		err := nodeInfoMap[node.Name].SetNode(node)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	for _, v := range nodeInfoMap {
// 		nodeInfos = append(nodeInfos, v)
// 	}

// 	return &testSharedLister{
// 		nodes:       nodes,
// 		nodeInfos:   nodeInfos,
// 		nodeInfoMap: nodeInfoMap,
// 	}
// }
