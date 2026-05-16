/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/

package cmd

import (
	"fmt"
	"os"
	"path"

	paths "github.com/dogukanmeral/scx-adapt/internal"
	"github.com/dogukanmeral/scx-adapt/internal/checks"
	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/spf13/cobra"
)

var listSchedulersCmd = &cobra.Command{
	Use:   "list-schedulers",
	Short: fmt.Sprintf("List schedulers in schedulers folder (%s)", paths.SCHEDULERSFOLDER),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			if os.Geteuid() != 0 {
				fmt.Println(MUST_RUN_AS_ROOT_MSG)
				os.Exit(1)
			}

			// Check if profiles directory exists
			if !helper.IsFileExist(paths.SCHEDULERSFOLDER) {
				fmt.Printf("ERROR: Schedulers folder '%s' does not exist.\n", paths.SCHEDULERSFOLDER)
				os.Exit(1)
			}

			// List external-loader schedulers
			fmt.Println("External-loader schedulers:")
			if helper.IsFileExist(paths.EXTERNALFOLDER) {

				externalFiles, err := os.ReadDir(paths.EXTERNALFOLDER)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// Iterate over external-loader schedulers and check
				for _, f := range externalFiles {
					if err := checks.CheckObj(path.Join(paths.EXTERNALFOLDER, f.Name())); err == nil {
						fmt.Printf("    %s\n", f.Name())
					}
				}
			}

			// List builtin-loader schedulers
			fmt.Println("Builtin-loader schedulers:")
			if helper.IsFileExist(paths.BUILTINFOLDER) {

				// Read entries in paths.BUILTIN directory
				builtinFiles, err := os.ReadDir(paths.BUILTINFOLDER)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// Iterate over builtin-loader schedulers and check
				for _, f := range builtinFiles {
					if checks.IsExecutableELF(path.Join(paths.BUILTINFOLDER, f.Name())) {
						fmt.Printf("    %s\n", f.Name())
					}
				}
			}
		default:
			fmt.Println(TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listSchedulersCmd)
}
