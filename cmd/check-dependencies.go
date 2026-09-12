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
	Long: `Check that all dependencies required by scx-adapt are available.

This verifies that the BPF filesystem is mounted and that the kernel has
sched_ext support (the '/sys/kernel/sched_ext' directory exists). If any
dependency is missing, it is reported.`,
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
