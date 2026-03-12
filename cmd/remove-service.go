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

// removeServiceCmd represents the removeService command
var removeServiceCmd = &cobra.Command{
	Use:   "remove-service",
	Short: fmt.Sprintf("Remove Systemd service file '%s' in '%s'", paths.SERVICEFILENAME, paths.SERVICESDIR),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}

		if os.Geteuid() != 0 {
			log.Fatalln("Must run as root")
		}

		// Check if .service file already exists.
		if !checks.IsFileExist(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME)) {
			log.Fatalf("Error: Service file does not exist at %s\n", path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))
		}

		// Remove service file.
		err := os.Remove(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Service file removed: %s\n", path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))
	},
}

func init() {
	rootCmd.AddCommand(removeServiceCmd)
}
