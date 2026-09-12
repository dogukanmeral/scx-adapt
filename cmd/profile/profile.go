package profile

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage scx-adapt profiles",
	}

	cmd.AddCommand(
		newAddCmd(),
		newLsCmd(),
		newRmCmd(),
		newCheckCmd(),
		newStartCmd(),
	)

	return cmd
}
