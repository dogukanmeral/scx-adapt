/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/dogukanmeral/scx-adapt/internal/msg"
	"github.com/dogukanmeral/scx-adapt/internal/paths"

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
			fmt.Println(msg.MISSING_ARGS_MSG)
			os.Exit(1)
		case 1:
			profileFile = args[0]
		default:
			fmt.Println(msg.TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println(msg.MUST_RUN_AS_ROOT_MSG)
			os.Exit(1)
		}

		// Check if profile exists in the profiles directory
		if !helper.IsFileExist(path.Join(paths.PROFILESFOLDER, profileFile)) {
			fmt.Printf("ERROR: Profile configuration with filename '%s' does not exist at '%s'\n",
				profileFile, paths.PROFILESFOLDER)
			os.Exit(1)
		}

		// Remove profile file in the profiles directory
		if err := os.Remove(path.Join(paths.PROFILESFOLDER, profileFile)); err != nil {
			fmt.Printf("ERROR: Deleting profile '%s' in '%s': %s\n",
				profileFile, paths.PROFILESFOLDER, err)
			os.Exit(1)
		}

		fmt.Printf("Profile at '%s' removed.\n", path.Join(paths.PROFILESFOLDER, profileFile))
	},
}

func init() {
	rootCmd.AddCommand(removeProfileCmd)
}
