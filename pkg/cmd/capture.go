package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Prajwal-kp-18/kubectl-meshsync-snapshot/pkg/kube"
	"github.com/Prajwal-kp-18/kubectl-meshsync-snapshot/pkg/meshsync"
	"github.com/spf13/cobra"
)

// NewCaptureCommand creates a new command for capturing MeshSync snapshot
func NewCaptureCommand() *cobra.Command {
	opts := &CaptureOptions{}

	cmd := &cobra.Command{
		Use:   "capture",
		Short: "Capture cluster state using MeshSync",
		Long:  `Capture the state of Kubernetes resources in the cluster using MeshSync.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCapture(opts)
		},
	}

	// Add flags specific to capture command
	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "meshery", "Namespace where MeshSync is deployed")
	cmd.Flags().StringVarP(&opts.OutputFile, "output", "o", "meshsync-snapshot.yaml", "Output file for snapshot")
	cmd.Flags().StringVarP(&opts.Format, "format", "f", "yaml", "Output format (yaml or json)")
	cmd.Flags().DurationVarP(&opts.Timeout, "timeout", "t", 60*time.Second, "Timeout for capture operation")
	cmd.Flags().BoolVarP(&opts.AllNamespaces, "all-namespaces", "A", false, "Capture resources from all namespaces")

	return cmd
}

// CaptureOptions contains options for capture command
type CaptureOptions struct {
	Namespace     string
	OutputFile    string
	Format        string
	Timeout       time.Duration
	AllNamespaces bool
}

// runCapture captures cluster state using MeshSync
func runCapture(opts *CaptureOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Create Kubernetes client
	client, err := kube.NewClient()
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}

	// Validate MeshSync is running
	if err := meshsync.Validate(ctx, client, opts.Namespace); err != nil {
		return fmt.Errorf("MeshSync validation failed: %w", err)
	}

	// Capture snapshot
	snapshot, err := meshsync.CaptureSnapshot(ctx, client, meshsync.CaptureOptions{
		Namespace:     opts.Namespace,
		AllNamespaces: opts.AllNamespaces,
	})
	if err != nil {
		return fmt.Errorf("failed to capture snapshot: %w", err)
	}

	// Save snapshot to file
	if err := meshsync.SaveSnapshot(snapshot, opts.OutputFile, opts.Format); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	fmt.Printf("Snapshot captured successfully and saved to %s\n", opts.OutputFile)
	return nil
} 