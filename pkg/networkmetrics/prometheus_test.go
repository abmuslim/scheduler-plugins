package networkmetrics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	defaultResultValue = 1441178050902
	defaultInterface   = "ens192"
	metricResult       = `
	{
        "metric": {
          "app": "prometheus",
          "app_kubernetes_io_managed_by": "Helm",
          "chart": "prometheus-13.6.0",
          "component": "node-exporter",
          "device": "ens192",
          "heritage": "Helm",
          "instance": "10.43.0.254:9100",
          "job": "kubernetes-service-endpoints",
          "kubernetes_name": "prometheus-1616380099-node-exporter",
          "kubernetes_namespace": "monitor",
          "kubernetes_node": "node05",
          "release": "prometheus-1616380099"
        },
        "value": [
          1618495272.297,
          "%d"
        ]
      }
	`
	successResponseTemplate = `
	{
		"status": "success",
		"data": {
		  "resultType": "vector",
		  "result": [%s]
		}
	  }
`
	errorResponse = `
{
	"status": "error",
	"errorType": "bad_data",
	"error": "invalid parameter \"query\": 1:90: parse error: bad duration syntax: \"5\""
} 
`
)

// Implements http.Handler
type MockPrometheus struct {
	t *testing.T
	*http.ServeMux
	requestsToError    int
	responseStatusCode int
	query              []byte
}

func NewMockPrometheusWithError(t *testing.T, requestsToError int, responseStatusCode int) MockPrometheus {
	return MockPrometheus{
		t:                  t,
		ServeMux:           http.NewServeMux(),
		requestsToError:    requestsToError,
		responseStatusCode: responseStatusCode,
	}
}

func NewMockPrometheus(t *testing.T) MockPrometheus {
	return NewMockPrometheusWithError(t, -1, http.StatusOK)
}

func (p *MockPrometheus) MockQueryRequestHandler(w http.ResponseWriter, _ *http.Request) {
	if p.shouldReturnError() {
		w.WriteHeader(p.responseStatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(p.query)
	if err != nil {
		p.t.Error("Error writing query:", err)
	}
}

func (p *MockPrometheus) shouldReturnError() bool {
	if p.requestsToError == 1 {
		return true
	}

	p.requestsToError--
	return false
}

func TestPrometheus(t *testing.T) {
	singleMetricResult := fmt.Sprintf(metricResult, defaultResultValue)
	multipleMetricsResult := fmt.Sprintf("%s,%s", singleMetricResult, singleMetricResult)
	testCases := []struct {
		name            string
		queryResult     string
		isErrorExpected bool
		requestsToError int
		statusCode      int
	}{
		{
			name:        "success scenario",
			queryResult: fmt.Sprintf(fmt.Sprintf(successResponseTemplate, metricResult), defaultResultValue),
		},
		{
			name:            "success response with multiple metrics",
			queryResult:     fmt.Sprintf(successResponseTemplate, multipleMetricsResult),
			isErrorExpected: true,
		},
		{
			name:            "query failure",
			queryResult:     errorResponse,
			isErrorExpected: true,
			requestsToError: 1,
			statusCode:      404,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPrometheus := NewMockPrometheus(t)
			if tc.isErrorExpected {
				mockPrometheus = NewMockPrometheusWithError(t, tc.requestsToError, tc.statusCode)
			}
			mockPrometheus.HandleFunc("/api/v1/query", mockPrometheus.MockQueryRequestHandler)
			mockPrometheus.query = []byte(tc.queryResult)
			mockServer := httptest.NewServer(mockPrometheus)
			defer mockServer.Close()

			prom := NewPrometheus(mockServer.URL, defaultInterface, time.Duration(5)*time.Minute)

			nodeMeasure, err := prom.getNodeBandwidthMeasure("node05")
			if tc.isErrorExpected {
				if err == nil {
					t.Error("error is expected executing prometheus query")
				}
			} else if int64(nodeMeasure.Value) != defaultResultValue {
				t.Errorf("expected response value is %d, got %s", defaultResultValue, nodeMeasure.Value)
			}
		})
	}

}
