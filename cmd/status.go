/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dogukanmeral/scx-adapt/internal/helper"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print currently running sched_ext scheduler.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			c, err := helper.CurrentScx()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Current scheduler: %s\n", c)
		default:
			fmt.Println(TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
