// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

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
