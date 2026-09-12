// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/dogukanmeral/scx-adapt/internal/msg"
	"github.com/dogukanmeral/scx-adapt/internal/paths"
	"github.com/spf13/cobra"
)

func newResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset the systemd service file to its default state",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Println(msg.TOO_MANY_ARGS_MSG)
				os.Exit(1)
			}

			if os.Geteuid() != 0 {
				fmt.Println(msg.MUST_RUN_AS_ROOT_MSG)
				os.Exit(1)
			}

			servicePath := path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME)
			dropInDir := path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME+".d")

			// Write default service file.
			if err := os.WriteFile(servicePath, []byte(SERVICEFILE), 0700); err != nil {
				fmt.Printf("ERROR: Writing service file '%s': %s\n", servicePath, err)
				os.Exit(1)
			}

			// Remove drop-in overrides created by 'systemctl edit'.
			if helper.IsFileExist(dropInDir) {
				if err := os.RemoveAll(dropInDir); err != nil {
					fmt.Printf("ERROR: Removing drop-in directory '%s': %s\n", dropInDir, err)
					os.Exit(1)
				}
			}

			fmt.Printf("Service file reset: %s\n", servicePath)

			reloadCmd := exec.Command("systemctl", "daemon-reload")
			if err := reloadCmd.Run(); err != nil {
				fmt.Printf("ERROR: Reloading daemons: %s\n", err)
				os.Exit(1)
			}
		},
	}
}
