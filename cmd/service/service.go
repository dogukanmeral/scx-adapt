package service

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the systemd service file",
	}

	cmd.AddCommand(
		newInstallCmd(),
		newRmCmd(),
	)

	return cmd
}
