/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
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
			fmt.Println(TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println(MUST_RUN_AS_ROOT_MSG)
			os.Exit(1)
		}

		if err := checks.CheckDependencies(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Check if .service file already exists.
		if checks.IsFileExist(path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME)) {
			fmt.Printf("Error: Service file already exists at %s\n", path.Join(paths.SERVICESDIR,
				paths.SERVICEFILENAME))
			os.Exit(1)
		}

		// Write service file.
		if err := os.WriteFile(path.Join(paths.SERVICESDIR,
			paths.SERVICEFILENAME), []byte(SERVICEFILE), 0700); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Service file added: %s\n", path.Join(paths.SERVICESDIR, paths.SERVICEFILENAME))

		reloadCmd := exec.Command("systemctl", "daemon-reload")

		if err := reloadCmd.Run(); err != nil {
			fmt.Printf("Error: Reloading daemons: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installServiceCmd)
}
