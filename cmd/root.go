/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
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
	Use:     "scx-adapt",
	Short:   "Adaptive and automated scheduler swtiching for sched_ext",
	Long:    ``,
	Version: "0.3.1",
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
