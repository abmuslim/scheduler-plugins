package pidcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"sigs.k8s.io/scheduler-plugins/pkg/apis/config"
)

const Name = "PIDController"
const MaxNodeScore = framework.MaxNodeScore

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
		return 0, framework.AsStatus(fmt.Errorf("error fetching score: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, framework.AsStatus(fmt.Errorf("error reading response: %v", err))
	}

	var result map[string]int
	if err = json.Unmarshal(body, &result); err != nil {
		return 0, framework.AsStatus(fmt.Errorf("error parsing JSON: %v", err))
	}

	score, ok := result["score"]
	if !ok {
		return 0, framework.AsStatus(fmt.Errorf("score not found in the response"))
	}

	return int64(score), nil
}

func (p *PIDController) ScoreExtensions() framework.ScoreExtensions {
	return p
}

func (p *PIDController) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	highest := int64(0)
	for _, score := range scores {
		if score.Score > highest {
			highest = score.Score
		}
	}

	if highest == 0 {
		return nil
	}

	for i := range scores {
		scores[i].Score = scores[i].Score * MaxNodeScore / highest
	}

	return nil
}

func New(ctx context.Context, obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	args, ok := obj.(*config.PIDControllerArgs)
	if !ok {
		return nil, fmt.Errorf("want args to be of type PIDControllerArgs, got %T", obj)
	}

	// Dereference pointer fields with proper nil checks to prevent panics
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
