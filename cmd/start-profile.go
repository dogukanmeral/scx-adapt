/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	paths "github.com/dogukanmeral/scx-adapt/internal"
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
			log.Fatalln("Missing arguments. scx-adapt --help to see usage")
		case 1:
			filepath = args[0]
		default:
			log.Fatalln("Too many arguments. scx-adapt --help to see usage")
		}

		if os.Geteuid() != 0 {
			log.Fatalln("Must run as root")
		}

		if err := checks.CheckDependencies(); err != nil {
			log.Fatalln(err)
		}

		// Check if lock exists (profiler already running)
		if checks.IsFileExist(paths.LOCKFILEPATH) {
			log.Fatalf("Error: Another scx-adapt profile already running. (%s)\n", paths.LOCKFILEPATH)
		}

		// Create DATAFOLDER folder if not exist
		if err := helper.CreateDirIfNotExist(paths.DATAFOLDER); err != nil {
			log.Fatalln(err)
		}

		// Create lock file
		if err := helper.CreateLock(); err != nil {
			log.Fatalln(err)
		}

		// If profile exists in PROFILESFOLDER with that name, use it
		if checks.IsFileExist(path.Join(paths.PROFILESFOLDER, filepath)) {
			filepath = path.Join(paths.PROFILESFOLDER, filepath)
		}

		yamlData, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatalf("Error occured while reading file '%s': %s\n", filepath, err)
		}

		conf, err := helper.YamlToConfig(yamlData)
		if err != nil {
			log.Fatalln(err)
		}

		// Interrupt handling
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		stop := make(chan bool, 1)
		errmsg := make(chan error, 1)
		schedChanged := make(chan helper.Scheduler, 1)

		go conf.Run(schedChanged, errmsg)

	STOPERROR:
		for {
			select {
			case err := <-errmsg:
				fmt.Println(err)

				if err := os.Remove(paths.LOCKFILEPATH); err != nil { // Remove the lock
					log.Fatalln("\nError: Removing lock file at 'scx-adapt.lock' failed.")
				}

			case sched := <-schedChanged:
				if checks.IsSchedExtActive() {
					stop <- true
					if err := <-errmsg; err != nil {
						errmsg <- err
						goto STOPERROR
					}
				}

				if sched.Path != "" {
					go sched.Run(stop, errmsg)
				}

			case <-interrupt:
				if checks.IsSchedExtActive() {
					stop <- true
					if err := <-errmsg; err != nil {
						errmsg <- err
						goto STOPERROR
					}
				}

				if err := os.Remove(paths.LOCKFILEPATH); err != nil { // Remove the lock
					fmt.Println("\nError: Removing lock file at 'scx-adapt.lock' failed.")
				}

				os.Exit(0)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startProfileCmd)
}
