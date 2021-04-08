package networkmetrics

import (
	"net/http"
	"testing"
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

//func TestPrometheus(t *testing.T) {
//	//mockPrometheus := NewMockPrometheus(t)
//	//mockPrometheus.HandleFunc("/", mockPrometheus.MockQueryRequestHandler)
//
//	prom := NewPrometheus("http://localhost:9090/")
//
//	prom.score("node1")
//}
