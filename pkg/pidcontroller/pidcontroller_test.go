package pidcontroller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"sigs.k8s.io/scheduler-plugins/pkg/apis/config"
)

func TestPIDController_Score(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		wantScore      int64
		expectError    bool
	}{
		{
			name:           "valid score",
			serverResponse: `{"score": 80}`,
			wantScore:      80,
			expectError:    false,
		},
		{
			name:           "missing score field",
			serverResponse: `{"value": 100}`,
			wantScore:      0,
			expectError:    true,
		},
		{
			name:           "malformed JSON",
			serverResponse: `{"score": 80`,
			wantScore:      0,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, tc.serverResponse)
			}))
			defer server.Close()

			args := &config.PIDControllerArgs{
				EndpointURL: stringPointer(server.URL),
			}

			plugin, err := New(context.Background(), args, nil)
			if err != nil {
				t.Fatalf("Failed to create plugin: %v", err)
			}

			pidController, ok := plugin.(*PIDController)
			if !ok {
				t.Fatalf("Plugin is not of type *PIDController")
			}

			score, status := pidController.Score(context.Background(), nil, nil, "")
			if tc.expectError && status.IsSuccess() {
				t.Errorf("Expected error, got success")
			}
			if !tc.expectError && !status.IsSuccess() {
				t.Errorf("Expected success, got error: %v", status.Message())
			}
			if score != tc.wantScore {
				t.Errorf("Expected score %d, got %d", tc.wantScore, score)
			}
		})
	}
}

func stringPointer(s string) *string {
	return &s
}
