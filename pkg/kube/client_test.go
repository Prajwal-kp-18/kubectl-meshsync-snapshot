package kube

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Save original KUBECONFIG environment variable
	originalKubeconfig := os.Getenv("KUBECONFIG")
	defer os.Setenv("KUBECONFIG", originalKubeconfig)

	tests := []struct {
		name        string
		kubeconfig  string
		expectError bool
	}{
		{
			name:        "Empty KUBECONFIG",
			kubeconfig:  "",
			expectError: false, // Should try to use default location
		},
		{
			name:        "Non-existent KUBECONFIG",
			kubeconfig:  "/non/existent/path/to/kubeconfig",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the KUBECONFIG environment variable for this test
			os.Setenv("KUBECONFIG", tt.kubeconfig)

			_, err := NewClient()
			if (err != nil) != tt.expectError {
				t.Errorf("NewClient() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
} 