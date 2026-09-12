// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package profile

import (
	"fmt"
	"os"

	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/dogukanmeral/scx-adapt/internal/msg"
	"github.com/spf13/cobra"
)

func newCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-profile <profile-path>",
		Short: "Check if profile file (in YAML format) is valid",
		Run: func(cmd *cobra.Command, args []string) {
			var profilePath string

			switch len(args) {
			case 0:
				fmt.Println(msg.MISSING_ARGS_MSG)
				os.Exit(1)
			case 1:
				profilePath = args[0]
			default:
				fmt.Println("Too many arguments. scx-adapt --help to see usage")
				os.Exit(1)
			}

			// Read file
			profileData, err := os.ReadFile(profilePath)
			if err != nil {
				fmt.Printf("ERROR: Reading file '%s': %s\n", profilePath, err)
				os.Exit(1)
			}

			// Check YAML configuration (discard "Config")
			if _, err = helper.YamlToConfig(profileData); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Valid config: %s\n", profilePath)
		},
	}
}
