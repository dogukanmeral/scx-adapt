/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/

package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	paths "github.com/dogukanmeral/scx-adapt/internal"
	"github.com/dogukanmeral/scx-adapt/internal/checks"
	"github.com/spf13/cobra"
)

var removeSchedulerCmd = &cobra.Command{
	Use:   "remove-scheduler <scheduler_filename>",
	Short: fmt.Sprintf("Remove scheduler from schedulers folder (%s)", paths.SCHEDULERSFOLDER),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var schedulerFile string

		switch len(args) {
		case 0:
			log.Fatalln("Missing arguments. scx-adapt --help to see usage")
		case 1:
			schedulerFile = args[0]
		default:
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}

		if os.Geteuid() != 0 {
			log.Fatalln("Must run as root")
		}

		// Check if scheduler exists in the schedulers directory
		if !checks.IsFileExist(path.Join(paths.SCHEDULERSFOLDER, schedulerFile)) {
			log.Fatalf("Scheduler with filename '%s' does not exist at '%s'\n", schedulerFile, paths.SCHEDULERSFOLDER)
		}

		// Remove scheduler file in the schedulers directory
		if err := os.Remove(path.Join(paths.SCHEDULERSFOLDER, schedulerFile)); err != nil {
			log.Fatalf("Error occured while deleting scheduler '%s' in '%s': %s\n", schedulerFile, paths.SCHEDULERSFOLDER, err)
		}

		fmt.Printf("Scheduler at '%s' removed.\n", path.Join(paths.SCHEDULERSFOLDER, schedulerFile))
	},
}

func init() {
	rootCmd.AddCommand(removeSchedulerCmd)
}
