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
	"github.com/dogukanmeral/scx-adapt/internal/checks"
	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/spf13/cobra"
)

var addSchedulerType string

var addSchedulerCmd = &cobra.Command{
	Use:   "add-scheduler [flags] <scheduler_path>",
	Short: fmt.Sprintf("Add sched_ext scheduler to schedulers folder (%s)", paths.SCHEDULERSFOLDER),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var schedulerPath string

		switch len(args) {
		case 0:
			fmt.Println(MISSING_ARGS_MSG)
			os.Exit(1)
		case 1:
			schedulerPath = args[0]
		default:
			fmt.Println(TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println(MUST_RUN_AS_ROOT_MSG)
			os.Exit(1)
		}

		var subdir string

		switch addSchedulerType {
		case string(helper.External):
			if err := checks.CheckObj(schedulerPath); err != nil {
				fmt.Printf("ERROR: Checking object file: %s\n", err)
				os.Exit(1)
			}

			subdir = paths.EXTERNALFOLDER

		case string(helper.Builtin):
			if !checks.IsExecutableELF(schedulerPath) {
				fmt.Printf("ERROR: Not an executable ELF file: %s\n", schedulerPath)
				os.Exit(1)
			}

			subdir = paths.BUILTINFOLDER

		default:
			fmt.Printf("ERROR: Invalid scheduler type '%s'. Available scheduler loader types: %s, %s\n",
				addSchedulerType, helper.External, helper.Builtin)
			os.Exit(1)
		}

		schedulerData, err := os.ReadFile(schedulerPath)
		if err != nil {
			fmt.Printf("ERROR: Reading file '%s': %s\n", schedulerPath, err)
			os.Exit(1)
		}

		// Check if a scheduler exists with the same name in schedulers directory
		if helper.IsFileExist(path.Join(subdir, filepath.Base(schedulerPath))) {
			fmt.Printf("ERROR: Another scheduler with filename '%s' already exists at '%s'\n", filepath.Base(schedulerPath), subdir)
			os.Exit(1)
		}

		// Create /etc/scx-adapt/ directory if not exist
		if err := helper.CreateDirIfNotExist(paths.DATAFOLDER); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create schedulers directory if not exist
		if err := helper.CreateDirIfNotExist(paths.SCHEDULERSFOLDER); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create scheduler type directory if not exist
		if err := helper.CreateDirIfNotExist(subdir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Copy file to schedulers directory
		if err := os.WriteFile(path.Join(subdir, filepath.Base(schedulerPath)), schedulerData, 0700); err != nil {
			fmt.Printf("ERROR: Writing to file '%s': %s\n", path.Join(subdir, filepath.Base(schedulerPath)), err)
			os.Exit(1)
		} else {
			fmt.Printf("Scheduler added to '%s'\n", path.Join(subdir, filepath.Base(schedulerPath)))
		}
	},
}

func init() {
	rootCmd.AddCommand(addSchedulerCmd)

	addSchedulerCmd.Flags().StringVarP(
		&addSchedulerType,
		"loader",
		"l",
		"",
		"Scheduler loader type (external|builtin)",
	)
	addSchedulerCmd.MarkFlagRequired("type")
}
