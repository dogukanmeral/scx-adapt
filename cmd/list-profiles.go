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

	"github.com/spf13/cobra"
)

// listProfilesCmd represents the listProfiles command
var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List profile configurations in profiles folder ('/etc/scx-adapt/profiles' by default)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			if os.Geteuid() != 0 {
				fmt.Println("Must run as root")
				os.Exit(1)
			}

			// Check if profiles directory exists
			if !checks.IsFileExist(PROFILESFOLDER) {
				fmt.Printf("Error: Profiles folder '%s' does not exist.\n", PROFILESFOLDER)
				os.Exit(1)
			}

			files, err := os.ReadDir(PROFILESFOLDER)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, f := range files {
				fileData, err := os.ReadFile(path.Join(PROFILESFOLDER, f.Name()))
				if err != nil {
					fmt.Println(err)
					continue
				}

				_, err = helper.YamlToConfig(fileData)
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println(f.Name())
			}
		default:
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listProfilesCmd)
}
