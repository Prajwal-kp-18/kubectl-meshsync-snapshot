package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Prajwal-kp-18/kubectl-meshsync-snapshot/pkg/kube"
	"github.com/Prajwal-kp-18/kubectl-meshsync-snapshot/pkg/meshsync"
	"github.com/spf13/cobra"
)

// NewCleanupCommand creates a new command for cleaning up MeshSync resources
func NewCleanupCommand() *cobra.Command {
	opts := &CleanupOptions{}

	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup MeshSync resources",
		Long:  `Remove MeshSync resources deployed for snapshot capture.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCleanup(opts)
		},
	}

	// Add flags specific to cleanup command
	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "meshery", "Namespace where MeshSync is deployed")
	cmd.Flags().DurationVarP(&opts.Timeout, "timeout", "t", 60*time.Second, "Timeout for cleanup operation")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Force cleanup even if resources are still in use")

	return cmd
}

// CleanupOptions contains options for cleanup command
type CleanupOptions struct {
	Namespace string
	Timeout   time.Duration
	Force     bool
}

// runCleanup removes MeshSync resources from the cluster
func runCleanup(opts *CleanupOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Create Kubernetes client
	client, err := kube.NewClient()
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}

	// Cleanup MeshSync resources
	err = meshsync.Cleanup(ctx, client, meshsync.CleanupOptions{
		Namespace: opts.Namespace,
		Force:     opts.Force,
	})
	if err != nil {
		return fmt.Errorf("failed to cleanup MeshSync resources: %w", err)
	}

	fmt.Printf("MeshSync resources cleaned up successfully from namespace %s\n", opts.Namespace)
	return nil
} 