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

var listSchedulersCmd = &cobra.Command{
	Use:   "list-schedulers",
	Short: fmt.Sprintf("List schedulers in schedulers folder (%s)", paths.SCHEDULERSFOLDER),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			if os.Geteuid() != 0 {
				log.Fatalln("Must run as root")
			}

			// Check if profiles directory exists
			if !checks.IsFileExist(paths.SCHEDULERSFOLDER) {
				log.Fatalf("Error: Schedulers folder '%s' does not exist.\n", paths.SCHEDULERSFOLDER)
			}

			files, err := os.ReadDir(paths.SCHEDULERSFOLDER)

			if err != nil {
				log.Fatalln(err)
			}

			for _, f := range files {
				if err := checks.CheckObj(path.Join(paths.SCHEDULERSFOLDER, f.Name())); err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println(f.Name())
			}
		default:
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}
	},
}

func init() {
	rootCmd.AddCommand(listSchedulersCmd)
}
