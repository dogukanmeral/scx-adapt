/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"internal/checks"
	"internal/helper"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addProfileCmd = &cobra.Command{
	Use:   "add-profile",
	Short: "Add scx-adapt profile configuration to profiles folder ('/etc/scx-adapt/profiles' by default)",
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

		if os.Geteuid() != 0 {
			fmt.Println("Must run as root")
			os.Exit(1)
		}

		// Check if configuration is valid
		profileData, err := os.ReadFile(profilePath)
		if err != nil {
			fmt.Printf("Error occured while reading file '%s': %s\n", profilePath, err)
			os.Exit(1)
		}

		_, err = helper.YamlToConfig(profileData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Check if a profile exists with the same name in profiles directory
		if checks.IsFileExist(path.Join(PROFILESFOLDER, filepath.Base(profilePath))) {
			fmt.Printf("Another profile configuration with filename '%s' already exists at '%s'\n", filepath.Base(profilePath), PROFILESFOLDER)
			os.Exit(1)
		}

		// Create /etc/scx-adapt/ directory if not exist
		if err := helper.CreateDirIfNotExist(DATAFOLDER); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create profiles directory if not exist
		if err := helper.CreateDirIfNotExist(PROFILESFOLDER); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Copy file to profiles directory
		if err := os.WriteFile(path.Join(PROFILESFOLDER, filepath.Base(profilePath)), profileData, 0700); err != nil {
			fmt.Printf("Error occured while writing to file '%s': %s\n", path.Join(PROFILESFOLDER, filepath.Base(profilePath)), err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addProfileCmd)
}
