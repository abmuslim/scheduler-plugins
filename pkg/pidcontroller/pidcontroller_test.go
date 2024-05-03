package pidcontroller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPIDController_Score ensures that the PIDController plugin fetches and calculates node scores correctly.
func TestPIDController_Score(t *testing.T) {
	tests := []struct {
		name          string
		serverResp    string
		expectedScore int64
		expectedErr   bool
	}{
		{
			name:          "valid score",
			serverResp:    `{"score": 5}`,
			expectedScore: 5,
			expectedErr:   false,
		},
		{
			name:          "invalid JSON",
			serverResp:    `{"score":}`,
			expectedScore: 0,
			expectedErr:   true,
		},
		{
			name:          "no score field",
			serverResp:    `{"invalid": 5}`,
			expectedScore: 0,
			expectedErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup a mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, tt.serverResp)
			}))
			defer server.Close()

			// Setup the PIDController with the mock server's URL
			p := &PIDController{
				client:      &http.Client{},
				endpointURL: server.URL,
			}

			// Execute the Score function
			score, status := p.Score(context.Background(), nil, nil, "fake-node")

			// Validate the results
			if tt.expectedErr && status.IsSuccess() {
				t.Errorf("expected error, got success")
			}
			if !tt.expectedErr && !status.IsSuccess() {
				t.Errorf("expected success, got error: %v", status.Message())
			}
			if score != tt.expectedScore {
				t.Errorf("expected score %d, got %d", tt.expectedScore, score)
			}
		})
	}
}
