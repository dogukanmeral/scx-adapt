/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dogukanmeral/scx-adapt/internal/helper"

	"github.com/dogukanmeral/scx-adapt/internal/checks"

	"github.com/spf13/cobra"
)

// startProfileCmd represents the startProfile command
var startProfileCmd = &cobra.Command{
	Use:   "start-profile <profile_path>",
	Short: "Run scx-adapt with the profile configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var filepath string

		switch len(args) {
		case 0:
			fmt.Println("Missing arguments. scx-adapt --help to see usage")
			os.Exit(1)
		case 1:
			filepath = args[0]
		default:
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println("Must run as root")
			os.Exit(1)
		}

		// Check if lock exists (profiler already running)
		if checks.IsFileExist(LOCKFILEPATH) {
			fmt.Printf("Error: Another scx-adapt profile already running. (%s)\n", LOCKFILEPATH)
			os.Exit(1)
		}

		// Interrupt handling
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-ch
			fmt.Printf("\nStopping profile '%s'...\n", filepath)

			if err := os.Remove(LOCKFILEPATH); err != nil { // Remove the lock
				fmt.Println("\nError: Removing lock file at 'scx-adapt.lock' failed.")
			}

			// Pass if no scx is running
			if checks.IsScxRunning() {
				if err := helper.StopCurrScx(); err != nil {
					fmt.Printf("\nError occured while stopping currently running sched_ext scheduler: %s\n", err)
					os.Exit(1)
				}
			}

			os.Exit(0)
		}()

		// Create /etc/scx-adapt/ folder if not exist
		if err := helper.CreateDirIfNotExist(DATAFOLDER); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create lock file
		if _, err := os.Create(LOCKFILEPATH); err != nil {
			fmt.Printf("Error occured while creating lock file at '%s': %s\n", LOCKFILEPATH, err)
		}

		// If profile exists in PROFILESFOLDER with that name, use it
		if checks.IsFileExist(path.Join(PROFILESFOLDER, filepath)) {
			filepath = path.Join(PROFILESFOLDER, filepath)
		}

		err := helper.RunProfile(filepath)

		if err != nil {
			fmt.Printf("Error occured while starting profile '%s': %s\n", filepath, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startProfileCmd)
}
