package kube

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Save original KUBECONFIG environment variable
	originalKubeconfig := os.Getenv("KUBECONFIG")
	defer os.Setenv("KUBECONFIG", originalKubeconfig)

	// Create a temporary kubeconfig file for testing
	tmpKubeconfig, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		t.Fatalf("Failed to create temp kubeconfig: %v", err)
	}
	defer os.Remove(tmpKubeconfig.Name())
	
	// Write invalid but parseable kubeconfig content
	_, err = tmpKubeconfig.WriteString(`
apiVersion: v1
kind: Config
clusters: []
users: []
contexts: []
current-context: ""
`)
	if err != nil {
		t.Fatalf("Failed to write to temp kubeconfig: %v", err)
	}
	tmpKubeconfig.Close()

	tests := []struct {
		name        string
		kubeconfig  string
		expectError bool
	}{
		{
			name:        "Valid KUBECONFIG",
			kubeconfig:  tmpKubeconfig.Name(),
			expectError: false,
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