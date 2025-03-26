package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Prajwal-kp-18/kubectl-meshsync-snapshot/pkg/kube"
	"github.com/Prajwal-kp-18/kubectl-meshsync-snapshot/pkg/meshsync"
	"github.com/spf13/cobra"
)

// NewDeployCommand creates a new command for deploying MeshSync
func NewDeployCommand() *cobra.Command {
	opts := &DeployOptions{}

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy MeshSync temporarily to cluster",
		Long:  `Deploy MeshSync component to cluster to capture kubernetes resources state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(opts)
		},
	}

	// Add flags specific to deploy command
	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "meshery", "Namespace to deploy MeshSync")
	cmd.Flags().StringVarP(&opts.Version, "version", "v", "latest", "MeshSync version to deploy")
	cmd.Flags().DurationVarP(&opts.Timeout, "timeout", "t", 120*time.Second, "Timeout for deployment")

	return cmd
}

// DeployOptions contains options for deploy command
type DeployOptions struct {
	Namespace string
	Version   string
	Timeout   time.Duration
}

// runDeploy deploys MeshSync to the cluster
func runDeploy(opts *DeployOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Create Kubernetes client
	client, err := kube.NewClient()
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}

	// Deploy MeshSync
	err = meshsync.Deploy(ctx, client, meshsync.DeployOptions{
		Namespace: opts.Namespace,
		Version:   opts.Version,
	})
	if err != nil {
		return fmt.Errorf("failed to deploy MeshSync: %w", err)
	}

	fmt.Printf("MeshSync deployed successfully in namespace %s\n", opts.Namespace)
	return nil
} 