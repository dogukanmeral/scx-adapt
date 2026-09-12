// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package service

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the systemd service file",
		Long: `Manage the systemd service file for scx-adapt.

Install, remove, edit, or reset the systemd unit that runs scx-adapt
automatically.`,
	}

	cmd.AddCommand(
		newInstallCmd(),
		newRmCmd(),
		newEditCmd(),
		newResetCmd(),
	)

	return cmd
}
