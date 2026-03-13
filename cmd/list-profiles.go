/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	paths "github.com/dogukanmeral/scx-adapt/internal"
	"github.com/dogukanmeral/scx-adapt/internal/helper"

	"github.com/dogukanmeral/scx-adapt/internal/checks"

	"github.com/spf13/cobra"
)

// listProfilesCmd represents the listProfiles command
var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: fmt.Sprintf("List profile configurations in profiles folder (%s)", paths.PROFILESFOLDER),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			if os.Geteuid() != 0 {
				fmt.Println(MUST_RUN_AS_ROOT_MSG)
				os.Exit(1)
			}

			// Check if profiles directory exists
			if !checks.IsFileExist(paths.PROFILESFOLDER) {
				fmt.Printf("Error: Profiles folder '%s' does not exist.\n", paths.PROFILESFOLDER)
				os.Exit(1)
			}

			files, err := os.ReadDir(paths.PROFILESFOLDER)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, f := range files {
				fileData, err := os.ReadFile(path.Join(paths.PROFILESFOLDER, f.Name()))
				if err != nil {
					fmt.Println(err)
					continue
				}

				_, err = helper.YamlToConfig(fileData)
				if err != nil {
					fmt.Printf("Error: In profile '%s': %s\n", f.Name(), err)
					continue
				}

				fmt.Println(f.Name())
			}
		default:
			fmt.Println(TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listProfilesCmd)
}
