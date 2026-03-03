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
	"github.com/spf13/cobra"
)

var removeSchedulerCmd = &cobra.Command{
	Use:   "remove-scheduler <scheduler_filename>",
	Short: "Remove scheduler from schedulers folder ('/etc/scx-adapt/schedulers' by default)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var schedulerFile string

		switch len(args) {
		case 0:
			fmt.Println("Missing arguments. scx-adapt --help to see usage")
			os.Exit(1)
		case 1:
			schedulerFile = args[0]
		default:
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println("Must run as root")
			os.Exit(1)
		}

		// Check if scheduler exists in the schedulers directory
		if !checks.IsFileExist(path.Join(paths.SCHEDULERSFOLDER, schedulerFile)) {
			fmt.Printf("Scheduler with filename '%s' does not exist at '%s'\n", schedulerFile, paths.SCHEDULERSFOLDER)
			os.Exit(1)
		}

		// Remove scheduler file in the schedulers directory
		if err := os.Remove(path.Join(paths.SCHEDULERSFOLDER, schedulerFile)); err != nil {
			fmt.Printf("Error occured while deleting scheduler '%s' in '%s': %s\n", schedulerFile, paths.SCHEDULERSFOLDER, err)
			os.Exit(1)
		}

		fmt.Printf("Scheduler at '%s' removed.\n", path.Join(paths.SCHEDULERSFOLDER, schedulerFile))
	},
}

func init() {
	rootCmd.AddCommand(removeSchedulerCmd)
}
