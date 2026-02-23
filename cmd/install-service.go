/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/dogukanmeral/scx-adapt/internal/checks"

	"github.com/spf13/cobra"
)

const SERVICEFILE string = `
[Unit]
Description=scx-adapt daemon for profile at %I

[Service]
Type=exec
ExecStart=/usr/bin/scx-adapt start-profile  %i

[Install]
WantedBy=multi-user.target`

const SERVICEFILENAME string = "scx-adapt@.service"
const SERVICESDIR string = "/etc/systemd/system"

// installServiceCmd represents the installService command
var installServiceCmd = &cobra.Command{
	Use:   "install-service",
	Short: fmt.Sprintf("Add Systemd service file '%s' to '%s'", SERVICEFILENAME, SERVICESDIR),
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
		if checks.IsFileExist(path.Join(SERVICESDIR, SERVICEFILENAME)) {
			fmt.Printf("Error: Service file already exists at %s\n", path.Join(SERVICESDIR, SERVICEFILENAME))
			os.Exit(1)
		}

		// Write service file.
		err := os.WriteFile(path.Join(SERVICESDIR, SERVICEFILENAME), []byte(SERVICEFILE), 0700)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Service file added: %s\n", path.Join(SERVICESDIR, SERVICEFILENAME))

		reloadCmd := exec.Command("systemctl", "daemon-reload")
		err = reloadCmd.Run()

		if err != nil {
			fmt.Printf("Error occured while reloading daemons: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installServiceCmd)
}
