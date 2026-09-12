// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>
package cmd

import (
	"fmt"
	"os"

	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/dogukanmeral/scx-adapt/internal/msg"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log-sched <file_path>",
	Short: "Write sched_ext event tracing to file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var filepath string

		if os.Geteuid() != 0 {
			fmt.Println(msg.MUST_RUN_AS_ROOT_MSG)
			os.Exit(1)
		}

		switch len(args) {
		case 0:
			fmt.Println(msg.MISSING_ARGS_MSG)
			os.Exit(1)
		case 1:
			filepath = args[0]
			helper.TraceSchedExt(filepath)
		default:
			fmt.Println(msg.TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
