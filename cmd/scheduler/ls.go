// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package scheduler

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

func newLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List added schedulers",
		Long: `List schedulers in the schedulers folder.

Prints the filename of each scheduler, grouped by loader type (external or
builtin). Only valid schedulers are listed.

Requires root privileges.`,
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				if os.Geteuid() != 0 {
					fmt.Println(msg.MUST_RUN_AS_ROOT_MSG)
					os.Exit(1)
				}

				// Check if profiles directory exists
				if !helper.IsFileExist(paths.SCHEDULERSFOLDER) {
					fmt.Printf("ERROR: Schedulers folder '%s' does not exist.\n", paths.SCHEDULERSFOLDER)
					os.Exit(1)
				}

				// List external-loader schedulers
				fmt.Println("External-loader schedulers:")
				if helper.IsFileExist(paths.EXTERNALFOLDER) {

					externalFiles, err := os.ReadDir(paths.EXTERNALFOLDER)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Iterate over external-loader schedulers and check
					for _, f := range externalFiles {
						if err := checks.CheckObj(path.Join(paths.EXTERNALFOLDER, f.Name())); err == nil {
							fmt.Printf("    %s\n", f.Name())
						}
					}
				}

				// List builtin-loader schedulers
				fmt.Println("Builtin-loader schedulers:")
				if helper.IsFileExist(paths.BUILTINFOLDER) {

					// Read entries in paths.BUILTIN directory
					builtinFiles, err := os.ReadDir(paths.BUILTINFOLDER)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					// Iterate over builtin-loader schedulers and check
					for _, f := range builtinFiles {
						if checks.IsExecutableELF(path.Join(paths.BUILTINFOLDER, f.Name())) {
							fmt.Printf("    %s\n", f.Name())
						}
					}
				}
			default:
				fmt.Println(msg.TOO_MANY_ARGS_MSG)
				os.Exit(1)
			}
		},
	}
}
