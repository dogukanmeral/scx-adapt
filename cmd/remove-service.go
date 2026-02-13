/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"internal/checks"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// removeServiceCmd represents the removeService command
var removeServiceCmd = &cobra.Command{
	Use:   "remove-service",
	Short: fmt.Sprintf("Remove Systemd service file '%s' in '%s'", SERVICEFILENAME, SERVICESDIR),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println("Must run as root")
			os.Exit(1)
		}

		// Check if .service file already exists.
		if !checks.IsFileExist(path.Join(SERVICESDIR, SERVICEFILENAME)) {
			fmt.Printf("Error: Service file does not exist at %s\n", path.Join(SERVICESDIR, SERVICEFILENAME))
			os.Exit(1)
		}

		// Remove service file.
		err := os.Remove(path.Join(SERVICESDIR, SERVICEFILENAME))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Service file removed: %s\n", path.Join(SERVICESDIR, SERVICEFILENAME))
	},
}

func init() {
	rootCmd.AddCommand(removeServiceCmd)
}
