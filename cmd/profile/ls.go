// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package profile

import (
	"fmt"
	"os"
	"path"

	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/dogukanmeral/scx-adapt/internal/msg"
	"github.com/dogukanmeral/scx-adapt/internal/paths"
	"github.com/spf13/cobra"
)

func newLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List profiles",
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				if os.Geteuid() != 0 {
					fmt.Println(msg.MUST_RUN_AS_ROOT_MSG)
					os.Exit(1)
				}

				// Check if profiles directory exists
				if !helper.IsFileExist(paths.PROFILESFOLDER) {
					fmt.Printf("ERROR: Profiles folder '%s' does not exist.\n", paths.PROFILESFOLDER)
					os.Exit(1)
				}

				// Read entries in paths.PROFILESFOLDER
				files, err := os.ReadDir(paths.PROFILESFOLDER)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// Iterate over all profiles and check if YAML structure is valid
				for _, f := range files {
					fileData, err := os.ReadFile(path.Join(paths.PROFILESFOLDER, f.Name()))
					if err != nil {
						fmt.Println(err)
						continue
					}

					_, err = helper.YamlToConfig(fileData)
					if err != nil {
						fmt.Printf("ERROR: In profile '%s': %s\n", f.Name(), err)
						continue
					}

					fmt.Println(f.Name())
				}
			default:
				fmt.Println(msg.TOO_MANY_ARGS_MSG)
				os.Exit(1)
			}
		},
	}
}
