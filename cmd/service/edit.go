// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package service

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dogukanmeral/scx-adapt/internal/msg"
	"github.com/dogukanmeral/scx-adapt/internal/paths"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit the systemd service file",
		Long: `Edit the systemd service file for scx-adapt.

Opens the unit file in the editor via 'systemctl edit --full'.

Requires root privileges.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Println(msg.TOO_MANY_ARGS_MSG)
				os.Exit(1)
			}

			if os.Geteuid() != 0 {
				fmt.Println(msg.MUST_RUN_AS_ROOT_MSG)
				os.Exit(1)
			}

			execCmd := exec.Command("systemctl", "edit", "--full", paths.SERVICEFILENAME)
			execCmd.Stdin = os.Stdin
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr

			if err := execCmd.Run(); err != nil {
				fmt.Printf("ERROR: Editing service file: %s\n", err)
				os.Exit(1)
			}
		},
	}
}
