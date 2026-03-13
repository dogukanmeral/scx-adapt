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

// removeServiceCmd represents the removeService command
var removeServiceCmd = &cobra.Command{
	Use:   "remove-service",
	Short: fmt.Sprintf("Remove Systemd service file '%s' in '%s'", paths.SERVICEFILENAME, paths.SERVICESDIR),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println(TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println(MUST_RUN_AS_ROOT_MSG)
			os.Exit(1)
		}

		// Check if .service file already exists.
		if !checks.IsFileExist(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME)) {
			fmt.Printf("Error: Service file does not exist at %s\n",
				path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))
			os.Exit(1)
		}

		// Remove service file.
		if err := os.Remove(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Service file removed: %s\n", path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))
	},
}

func init() {
	rootCmd.AddCommand(removeServiceCmd)
}
