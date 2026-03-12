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

var addSchedulerCmd = &cobra.Command{
	Use:   "add-scheduler <scheduler_path>",
	Short: fmt.Sprintf("Add sched_ext scheduler object file to schedulers folder (%s)", paths.SCHEDULERSFOLDER),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var schedulerPath string

		switch len(args) {
		case 0:
			log.Fatalln("Missing arguments. scx-adapt --help to see usage")
		case 1:
			schedulerPath = args[0]
		default:
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}

		if os.Geteuid() != 0 {
			log.Fatalln("Must run as root")
		}

		// Check obj
		if err := checks.CheckObj(schedulerPath); err != nil {
			log.Fatalf("Error: Checking object file: %s\n", err)
		}

		schedulerData, err := os.ReadFile(schedulerPath)
		if err != nil {
			log.Fatalf("Error: Reading file '%s': %s\n", schedulerPath, err)
		}

		// Check if a scheduler exists with the same name in schedulers directory
		if checks.IsFileExist(path.Join(paths.SCHEDULERSFOLDER, filepath.Base(schedulerPath))) {
			log.Fatalf("Another scheduler with filename '%s' already exists at '%s'\n", filepath.Base(schedulerPath), paths.SCHEDULERSFOLDER)
		}

		// Create /etc/scx-adapt/ directory if not exist
		if err := helper.CreateDirIfNotExist(paths.DATAFOLDER); err != nil {
			log.Fatalln(err)
		}

		// Create schedulers directory if not exist
		if err := helper.CreateDirIfNotExist(paths.SCHEDULERSFOLDER); err != nil {
			log.Fatalln(err)
		}

		// Copy file to profiles directory
		if err := os.WriteFile(path.Join(paths.SCHEDULERSFOLDER, filepath.Base(schedulerPath)), schedulerData, 0700); err != nil {
			log.Fatalf("Error: Writing to file '%s': %s\n", path.Join(paths.SCHEDULERSFOLDER, filepath.Base(schedulerPath)), err)
		} else {
			fmt.Printf("Profile added to '%s'\n", path.Join(paths.SCHEDULERSFOLDER, filepath.Base(schedulerPath)))
		}
	},
}

func init() {
	rootCmd.AddCommand(addSchedulerCmd)
}
