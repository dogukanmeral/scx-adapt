// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package service

import (
	"fmt"
	"os"
	"path"

	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/dogukanmeral/scx-adapt/internal/msg"
	"github.com/dogukanmeral/scx-adapt/internal/paths"
	"github.com/spf13/cobra"
)

func newRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm",
		Short: "Remove Systemd service file",
		Long: `Remove the systemd service file for scx-adapt.

Deletes the scx-adapt systemd unit from the systemd directory.

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

			// Check if .service file already exists.
			if !helper.IsFileExist(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME)) {
				fmt.Printf("ERROR: Service file does not exist at %s\n",
					path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))
				os.Exit(1)
			}

			// Remove service file.
			if err := os.Remove(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME)); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Service file removed: %s\n", path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))
		},
	}
}
