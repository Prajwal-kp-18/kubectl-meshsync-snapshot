package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Prajwal-kp-18/kubectl-meshsync-snapshot/pkg/meshery"
	"github.com/spf13/cobra"
)

// NewImportCommand creates a new command for importing snapshot to Meshery
func NewImportCommand() *cobra.Command {
	opts := &ImportOptions{}

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import snapshot to Meshery",
		Long:  `Import captured snapshot to Meshery via API or file output.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImport(opts)
		},
	}

	// Add flags specific to import command
	cmd.Flags().StringVarP(&opts.MesheryURL, "url", "u", "http://localhost:9081", "Meshery server URL")
	cmd.Flags().StringVarP(&opts.Token, "token", "t", "", "Meshery authentication token")
	cmd.Flags().StringVarP(&opts.InputFile, "input", "i", "meshsync-snapshot.yaml", "Input snapshot file path")
	cmd.Flags().DurationVarP(&opts.Timeout, "timeout", "", 30*time.Second, "Timeout for import operation")

	return cmd
}

// ImportOptions contains options for import command
type ImportOptions struct {
	MesheryURL string
	Token      string
	InputFile  string
	Timeout    time.Duration
}

// runImport imports snapshot to Meshery
func runImport(opts *ImportOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Create Meshery client
	client, err := meshery.NewClient(opts.MesheryURL, opts.Token)
	if err != nil {
		return fmt.Errorf("error creating Meshery client: %w", err)
	}

	// Import snapshot
	err = client.ImportSnapshot(ctx, opts.InputFile)
	if err != nil {
		return fmt.Errorf("failed to import snapshot: %w", err)
	}

	fmt.Printf("Snapshot imported successfully to Meshery server %s\n", opts.MesheryURL)
	return nil
} 