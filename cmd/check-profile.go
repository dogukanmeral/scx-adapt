/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/dogukanmeral/scx-adapt/internal/helper"

	"github.com/spf13/cobra"
)

// checkProfileCmd represents the checkProfile command
var checkProfileCmd = &cobra.Command{
	Use:   "check-profile <profile_path>",
	Short: "Check if profile file in YAML format passed from STDIN is valid",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println("Error reading from stdin: ", err)
				os.Exit(1)
			}

			if _, err := helper.YamlToConfig(data); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Valid configuration.")
		} else {
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkProfileCmd)
}
