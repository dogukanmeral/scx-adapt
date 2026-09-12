// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>
package cmd

import (
	"fmt"
	"os"

	"github.com/dogukanmeral/scx-adapt/cmd/profile"
	"github.com/dogukanmeral/scx-adapt/cmd/scheduler"
	"github.com/dogukanmeral/scx-adapt/cmd/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scx-adapt",
	Short: "Adaptive and automated scheduler swtiching for sched_ext",
	Long: `scx-adapt is an adaptive and automated scheduler policy manager for sched_ext.

It manages profiles (sets of sched_ext schedulers with selection criteria),
schedulers, and the systemd service used to run them automatically.`,
	Version: "0.4.0",
}

func init() {
	rootCmd.AddCommand(
		profile.NewCmd(),
		scheduler.NewCmd(),
		service.NewCmd(),
	)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
