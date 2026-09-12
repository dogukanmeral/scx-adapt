// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

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
