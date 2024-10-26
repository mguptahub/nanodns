package dns

import (
	"testing"
)

func TestResolveServiceIP(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		wantErr     bool
	}{
		{
			name:        "invalid service name",
			serviceName: "nonexistent-service",
			wantErr:     true,
		},
		// Note: Can't easily test successful resolution without a running Docker environment
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := ResolveServiceIP(tt.serviceName)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ResolveServiceIP() expected error for service %q, got nil", tt.serviceName)
				}
				return
			}
			if err != nil {
				t.Errorf("ResolveServiceIP() error = %v, want nil", err)
			}
			if !tt.wantErr && ip == "" {
				t.Error("ResolveServiceIP() returned empty IP for valid service")
			}
		})
	}
}
