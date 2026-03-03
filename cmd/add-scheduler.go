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

var addSchedulerCmd = &cobra.Command{
	Use:   "add-scheduler <scheduler_path>",
	Short: fmt.Sprintf("Add sched_ext scheduler object file to schedulers folder (%s)", paths.SCHEDULERSFOLDER),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var schedulerPath string

		switch len(args) {
		case 0:
			fmt.Println("Missing arguments. scx-adapt --help to see usage")
			os.Exit(1)
		case 1:
			schedulerPath = args[0]
		default:
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println("Must run as root")
			os.Exit(1)
		}

		// Check obj
		if err := checks.CheckObj(schedulerPath); err != nil {
			fmt.Printf("Error occured while checking object file: %s\n", err)
			os.Exit(1)
		}

		schedulerData, err := os.ReadFile(schedulerPath)
		if err != nil {
			fmt.Printf("Error occured while reading file '%s': %s\n", schedulerPath, err)
			os.Exit(1)
		}

		// Check if a scheduler exists with the same name in schedulers directory
		if checks.IsFileExist(path.Join(paths.SCHEDULERSFOLDER, filepath.Base(schedulerPath))) {
			fmt.Printf("Another scheduler with filename '%s' already exists at '%s'\n", filepath.Base(schedulerPath), paths.SCHEDULERSFOLDER)
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

		// Copy file to profiles directory
		if err := os.WriteFile(path.Join(paths.SCHEDULERSFOLDER, filepath.Base(schedulerPath)), schedulerData, 0700); err != nil {
			fmt.Printf("Error occured while writing to file '%s': %s\n", path.Join(paths.SCHEDULERSFOLDER, filepath.Base(schedulerPath)), err)
			os.Exit(1)
		} else {
			fmt.Printf("Profile added to '%s'\n", path.Join(paths.SCHEDULERSFOLDER, filepath.Base(schedulerPath)))
		}
	},
}

func init() {
	rootCmd.AddCommand(addSchedulerCmd)
}
