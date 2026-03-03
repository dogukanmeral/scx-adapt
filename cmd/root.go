/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scx-adapt",
	Short: "Adaptive and automated scheduler policies for sched_ext",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
