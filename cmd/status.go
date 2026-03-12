/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"log"

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
				log.Fatalln(err)
			}

			fmt.Printf("Current scheduler: %s\n", c)
		default:
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
