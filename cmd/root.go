/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const MISSING_ARGS_MSG = "Missing arguments. scx-adapt --help to see usage"
const TOO_MANY_ARGS_MSG = "Too many arguments. scx-adapt --help to see usage"
const MUST_RUN_AS_ROOT_MSG = "Must run as root"
const INTERRUPT_MSG = "Interrupted... Exiting..."

var rootCmd = &cobra.Command{
	Use:   "scx-adapt",
	Short: "Adaptive and automated scheduler policies for sched_ext",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {}
