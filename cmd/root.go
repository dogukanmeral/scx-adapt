/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const PROFILESFOLDER string = "/etc/scx-adapt/profiles/"
const DATAFOLDER string = "/etc/scx-adapt/"
const LOCKFILEPATH string = "/etc/scx-adapt/scx-adapt.lock"

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
