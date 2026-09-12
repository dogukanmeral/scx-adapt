// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>
package cmd

import (
	"fmt"
	"os"

	"github.com/dogukanmeral/scx-adapt/internal/checks"
	"github.com/spf13/cobra"
)

// checkDependenciesCmd represents the checkDepencies command
var checkDependenciesCmd = &cobra.Command{
	Use:   "check-dependencies",
	Short: "Check dependencies of scx-adapt",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}

		checks.CheckBPFDependencies()
	},
}

func init() {
	rootCmd.AddCommand(checkDependenciesCmd)
}
