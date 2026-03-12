/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/dogukanmeral/scx-adapt/internal/helper"

	"github.com/spf13/cobra"
)

// checkProfileCmd represents the checkProfile command
var checkProfileCmd = &cobra.Command{
	Use:   "check-profile <profile_path>",
	Short: "Check if profile file in YAML format is valid",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var profilePath string

		switch len(args) {
		case 0:
			log.Fatalln("Missing arguments. scx-adapt --help to see usage")
		case 1:
			profilePath = args[0]
		default:
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}

		// Read file
		profileData, err := os.ReadFile(profilePath)
		if err != nil {
			log.Fatalf("Error: Reading file '%s': %s\n", profilePath, err)
		}

		// Check YAML configuration (discard "Config")
		_, err = helper.YamlToConfig(profileData)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Valid config: %s\n", profilePath)
	},
}

func init() {
	rootCmd.AddCommand(checkProfileCmd)
}
