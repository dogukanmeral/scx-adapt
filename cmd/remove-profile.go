/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"internal/checks"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var removeProfileCmd = &cobra.Command{
	Use:   "remove-profile",
	Short: "Remove profile configuration from profiles folder ('/etc/scx-adapt/profiles' by default)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var profileFile string

		switch len(args) {
		case 0:
			fmt.Println("Missing arguments. scx-adapt --help to see usage")
			os.Exit(1)
		case 1:
			profileFile = args[0]
		default:
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println("Must run as root")
			os.Exit(1)
		}

		// Check if profile exists in the profiles directory
		if !checks.IsFileExist(path.Join(PROFILESFOLDER, profileFile)) {
			fmt.Printf("Profile configuration with filename '%s' does not exist at '%s'\n", profileFile, PROFILESFOLDER)
			os.Exit(1)
		}

		// Remove profile file in the profiles directory
		if err := os.Remove(path.Join(PROFILESFOLDER, profileFile)); err != nil {
			fmt.Printf("Error occured while deleting profile '%s' in '%s': %s\n", profileFile, PROFILESFOLDER, err)
			os.Exit(1)
		}

		fmt.Printf("Profile at '%s' removed.\n", path.Join(PROFILESFOLDER, profileFile))
	},
}

func init() {
	rootCmd.AddCommand(removeProfileCmd)
}
