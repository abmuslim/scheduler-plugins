package pidcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"k8s.io/klog/v2" // Import klog for logging

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/scheduler-plugins/apis/config"
)

const Name = "PIDController"
const MaxNodeScore = 100 // Updated max score to 100 to reflect normalization

type PIDController struct {
	handle      framework.Handle
	client      *http.Client
	endpointURL string
}

func (p *PIDController) Name() string {
	return "PIDController"
}

func (p *PIDController) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	// Send GET request to the configured endpoint URL
	resp, err := p.client.Get(p.endpointURL)
	if err != nil {
		klog.Errorf("Failed to fetch score: %v", err)
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("error fetching score: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("Failed to read response body: %v", err)
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("error reading response: %v", err))
	}

	var result map[string]int
	if err = json.Unmarshal(body, &result); err != nil {
		klog.Errorf("Failed to parse JSON: %v", err)
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("error parsing JSON: %v", err))
	}

	score, ok := result["score"]
	if !ok {
		klog.Errorf("Score not found in the response")
		return 0, framework.NewStatus(framework.Error, "score not found in the response")
	}

	// Normalize the score from a 1-10 scale to a 1-100 scale
	normalizedScore := int64(score) * 10
	return normalizedScore, nil
}

func (p *PIDController) ScoreExtensions() framework.ScoreExtensions {
	return p
}

// The NormalizeScore function can be simplified or removed if normalization is handled in Score method
func (p *PIDController) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	return nil // Since normalization is now handled in the Score method directly
}

func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	args, ok := obj.(*config.PIDControllerArgs)
	if !ok {
		return nil, fmt.Errorf("want args to be of type PIDControllerArgs, got %T", obj)
	}

	var maxIdleConns, idleTimeoutSec, requestTimeoutSec int
	var endpointURL string

	if args.MaxIdleConnections != nil {
		maxIdleConns = *args.MaxIdleConnections
	}
	if args.IdleConnectionTimeoutSec != nil {
		idleTimeoutSec = *args.IdleConnectionTimeoutSec
	}
	if args.RequestTimeoutSec != nil {
		requestTimeoutSec = *args.RequestTimeoutSec
	}
	if args.EndpointURL != nil {
		endpointURL = *args.EndpointURL
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    maxIdleConns,
			IdleConnTimeout: time.Duration(idleTimeoutSec) * time.Second,
		},
		Timeout: time.Duration(requestTimeoutSec) * time.Second,
	}

	return &PIDController{
		handle:      handle,
		client:      client,
		endpointURL: endpointURL,
	}, nil
}
