package scheduler

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduler",
		Short: "Manage schedulers",
	}

	cmd.AddCommand(
		newAddCmd(),
		newLsCmd(),
		newRmCmd(),
	)

	return cmd
}
