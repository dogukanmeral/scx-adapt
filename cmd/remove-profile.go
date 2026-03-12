/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	paths "github.com/dogukanmeral/scx-adapt/internal"
	"github.com/dogukanmeral/scx-adapt/internal/checks"

	"github.com/spf13/cobra"
)

var removeProfileCmd = &cobra.Command{
	Use:   "remove-profile <profile_filename>",
	Short: fmt.Sprintf("Remove profile configuration from profiles folder (%s)", paths.PROFILESFOLDER),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var profileFile string

		switch len(args) {
		case 0:
			log.Fatalln("Missing arguments. scx-adapt --help to see usage")
		case 1:
			profileFile = args[0]
		default:
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}

		if os.Geteuid() != 0 {
			log.Fatalln("Must run as root")
		}

		// Check if profile exists in the profiles directory
		if !checks.IsFileExist(path.Join(paths.PROFILESFOLDER, profileFile)) {
			log.Fatalf("Profile configuration with filename '%s' does not exist at '%s'\n", profileFile, paths.PROFILESFOLDER)
		}

		// Remove profile file in the profiles directory
		if err := os.Remove(path.Join(paths.PROFILESFOLDER, profileFile)); err != nil {
			log.Fatalf("Error occured while deleting profile '%s' in '%s': %s\n", profileFile, paths.PROFILESFOLDER, err)
		}

		fmt.Printf("Profile at '%s' removed.\n", path.Join(paths.PROFILESFOLDER, profileFile))
	},
}

func init() {
	rootCmd.AddCommand(removeProfileCmd)
}
