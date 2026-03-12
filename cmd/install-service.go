/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	paths "github.com/dogukanmeral/scx-adapt/internal"
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

// installServiceCmd represents the installService command
var installServiceCmd = &cobra.Command{
	Use:   "install-service",
	Short: fmt.Sprintf("Add Systemd service file '%s' to '%s'", paths.SERVICEFILENAME, paths.SERVICESDIR),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}

		if os.Geteuid() != 0 {
			log.Fatalln("Must run as root")
		}

		if err := checks.CheckDependencies(); err != nil {
			log.Fatalln(err)
		}

		// Check if .service file already exists.
		if checks.IsFileExist(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME)) {
			log.Fatalf("Error: Service file already exists at %s\n", path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))
		}

		// Write service file.
		err := os.WriteFile(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME), []byte(SERVICEFILE), 0700)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Service file added: %s\n", path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))

		reloadCmd := exec.Command("systemctl", "daemon-reload")
		err = reloadCmd.Run()

		if err != nil {
			log.Fatalf("Error occured while reloading daemons: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(installServiceCmd)
}
