// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/dogukanmeral/scx-adapt/internal/checks"
	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/dogukanmeral/scx-adapt/internal/msg"
	"github.com/dogukanmeral/scx-adapt/internal/paths"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the live state and saved objects of scx-adapt",
	Long: `Print the live state and saved objects of scx-adapt.

Reports the number of builtin-loader and external-loader schedulers, the
number of profiles, the active sched_ext scheduler, the active profile, and
whether the systemd service file is installed.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			if os.Geteuid() != 0 {
				fmt.Println(msg.MUST_RUN_AS_ROOT_MSG)
				os.Exit(1)
			}

			fmt.Printf("Builtin-loader schedulers:  %d\n", countFiles(paths.BUILTINFOLDER))
			fmt.Printf("External-loader schedulers: %d\n", countFiles(paths.EXTERNALFOLDER))
			fmt.Printf("Profiles:                   %d\n", countFiles(paths.PROFILESFOLDER))

			// Current state of sched_ext
			if checks.IsSchedExtActive() {
				c, err := helper.CurrentScx()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("sched_ext:                  active (%s)\n", c)
			} else {
				fmt.Println("sched_ext:                  inactive")
			}

			// Current state of scx-adapt
			if helper.IsFileExist(paths.LOCKFILEPATH) {
				name, err := helper.ReadLock()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("scx-adapt:                  running (%s)\n", name)
			} else {
				fmt.Println("scx-adapt:                  not running")
			}

			// Service file installed or not
			if helper.IsFileExist(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME)) {
				fmt.Println("Service file:               installed")
			} else {
				fmt.Println("Service file:               not installed")
			}
		default:
			fmt.Println(msg.TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}
	},
}

// countFiles returns the number of entries in a directory (0 if it does not exist).
func countFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	return len(entries)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
