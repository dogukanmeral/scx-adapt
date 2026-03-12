/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"log"
	"os"

	"github.com/dogukanmeral/scx-adapt/internal/helper"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log-sched <file_path>",
	Short: "Write sched_ext event tracing to file.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var filepath string

		if os.Geteuid() != 0 {
			log.Fatalln("Must run as root")
		}

		switch len(args) {
		case 0:
			log.Fatalln("Missing arguments. scx-adapt --help to see usage")
		case 1:
			filepath = args[0]
			helper.TraceSchedExt(filepath)
		default:
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
