package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand creates the root command for meshsync-snapshot plugin
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meshsync-snapshot",
		Short: "A kubectl plugin for capturing MeshSync snapshots",
		Long: `meshsync-snapshot is a kubectl plugin that helps manage 
MeshSync snapshots for Meshery. It allows deploying MeshSync,
capturing cluster state, importing snapshots and cleaning up resources.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(NewDeployCommand())
	cmd.AddCommand(NewCaptureCommand())
	cmd.AddCommand(NewImportCommand())
	cmd.AddCommand(NewCleanupCommand())

	return cmd
} 