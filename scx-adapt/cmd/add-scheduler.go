/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"internal/helper"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add-scheduler",
	Short: "Add scheduler(s) (BPF bytecode) to scx-adapt",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			for _, scx := range args {
				err := helper.AddScx(scx, "/var/lib/scx-adapt/obj")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println("Scheduler added.")
			}
		} else {
			fmt.Println("Missing scheduler file(s) as argument. scx-adapt --help to see usage")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
