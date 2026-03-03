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

// checkProfileCmd represents the checkProfile command
var checkProfileCmd = &cobra.Command{
	Use:   "check-profile <profile_path>",
	Short: "Check if profile file in YAML format passed from STDIN is valid",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var profilePath string

		switch len(args) {
		case 0:
			fmt.Println("Missing arguments. scx-adapt --help to see usage")
			os.Exit(1)
		case 1:
			profilePath = args[0]
		default:
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}

		// Read file
		profileData, err := os.ReadFile(profilePath)
		if err != nil {
			fmt.Printf("Error occured while reading file '%s': %s\n", profilePath, err)
			os.Exit(1)
		}

		// Check YAML configuration (discard "Config")
		_, err = helper.YamlToConfig(profileData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Valid config: %s\n", profilePath)
	},
}

func init() {
	rootCmd.AddCommand(checkProfileCmd)
}
