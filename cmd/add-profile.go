/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	paths "github.com/dogukanmeral/scx-adapt/internal"
	"github.com/dogukanmeral/scx-adapt/internal/checks"
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
			log.Fatalln("Missing arguments. scx-adapt --help to see usage")
		case 1:
			profilePath = args[0]
		default:
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}

		if os.Geteuid() != 0 {
			log.Fatalln("Must run as root")
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

		// Check if a profile exists with the same name in profiles directory
		if checks.IsFileExist(path.Join(paths.PROFILESFOLDER, filepath.Base(profilePath))) {
			log.Fatalf("Another profile configuration with filename '%s' already exists at '%s'\n", filepath.Base(profilePath), paths.PROFILESFOLDER)
		}

		// Create /etc/scx-adapt/ directory if not exist
		if err := helper.CreateDirIfNotExist(paths.DATAFOLDER); err != nil {
			log.Fatalln(err)
		}

		// Create profiles directory if not exist
		if err := helper.CreateDirIfNotExist(paths.PROFILESFOLDER); err != nil {
			log.Fatalln(err)
		}

		// Copy file to profiles directory
		if err := os.WriteFile(path.Join(paths.PROFILESFOLDER, filepath.Base(profilePath)), profileData, 0700); err != nil {
			log.Fatalf("Error: Writing to file '%s': %s\n", path.Join(paths.PROFILESFOLDER, filepath.Base(profilePath)), err)
		} else {
			fmt.Printf("Profile added to '%s'\n", path.Join(paths.PROFILESFOLDER, filepath.Base(profilePath)))
		}
	},
}

func init() {
	rootCmd.AddCommand(addProfileCmd)
}
