/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	paths "github.com/dogukanmeral/scx-adapt/internal"
	"github.com/dogukanmeral/scx-adapt/internal/helper"

	"github.com/spf13/cobra"
)

var addProfileCmd = &cobra.Command{
	Use:   "add-profile <profile_path>",
	Short: fmt.Sprintf("Add scx-adapt profile configuration to profiles folder (%s)", paths.PROFILESFOLDER),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var profilePath string

		switch len(args) {
		case 0:
			fmt.Println(MISSING_ARGS_MSG)
			os.Exit(1)
		case 1:
			profilePath = args[0]
		default:
			fmt.Println(TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println(MUST_RUN_AS_ROOT_MSG)
			os.Exit(1)
		}

		// Read file
		profileData, err := os.ReadFile(profilePath)
		if err != nil {
			fmt.Printf("Error: Reading file '%s': %s\n", profilePath, err)
			os.Exit(1)
		}

		// Check YAML configuration (discard "Config")
		_, err = helper.YamlToConfig(profileData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Check if a profile exists with the same name in profiles directory
		if helper.IsFileExist(path.Join(paths.PROFILESFOLDER, filepath.Base(profilePath))) {
			fmt.Printf("Another profile configuration with filename '%s' already exists at '%s'\n", filepath.Base(profilePath), paths.PROFILESFOLDER)
			os.Exit(1)
		}

		// Create /etc/scx-adapt/ directory if not exist
		if err := helper.CreateDirIfNotExist(paths.DATAFOLDER); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create profiles directory if not exist
		if err := helper.CreateDirIfNotExist(paths.PROFILESFOLDER); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Copy file to profiles directory
		if err := os.WriteFile(path.Join(paths.PROFILESFOLDER, filepath.Base(profilePath)), profileData, 0700); err != nil {
			fmt.Printf("Error: Writing to file '%s': %s\n", path.Join(paths.PROFILESFOLDER, filepath.Base(profilePath)), err)
			os.Exit(1)
		} else {
			fmt.Printf("Profile added to '%s'\n", path.Join(paths.PROFILESFOLDER, filepath.Base(profilePath)))
		}
	},
}

func init() {
	rootCmd.AddCommand(addProfileCmd)
}
