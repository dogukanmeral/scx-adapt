// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package profile

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage scx-adapt profiles",
		Long: `Manage scx-adapt profiles.

Profiles are YAML configuration files that define a set of sched_ext
schedulers and the criteria used to select between them at runtime.`,
	}

	cmd.AddCommand(
		newAddCmd(),
		newLsCmd(),
		newRmCmd(),
		newEditCmd(),
		newCheckCmd(),
		newStartCmd(),
	)

	return cmd
}
