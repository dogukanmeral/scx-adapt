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
				fmt.Printf("Error: Schedulers folder '%s' does not exist.\n", paths.SCHEDULERSFOLDER)
				os.Exit(1)
			}

			// List kernelonly schedulers
			fmt.Println("Kernel-only schedulers:")
			if helper.IsFileExist(paths.KERNELONLYFOLDER) {

				kernelonlyFiles, err := os.ReadDir(paths.KERNELONLYFOLDER)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				for _, f := range kernelonlyFiles {
					if err := checks.CheckObj(path.Join(paths.KERNELONLYFOLDER, f.Name())); err == nil {
						fmt.Printf("    %s\n", f.Name())
					}
				}
			}

			// List userspace schedulers
			fmt.Println("Userspace schedulers:")
			if helper.IsFileExist(paths.USERSPACEFOLDER) {

				userspaceFiles, err := os.ReadDir(paths.USERSPACEFOLDER)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				for _, f := range userspaceFiles {
					if checks.IsExecutableELF(path.Join(paths.USERSPACEFOLDER, f.Name())) {
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
