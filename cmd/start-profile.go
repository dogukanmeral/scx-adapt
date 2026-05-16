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
			fmt.Println(MISSING_ARGS_MSG)
			os.Exit(1)
		case 1:
			filepath = args[0]
		default:
			fmt.Println(TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println(MUST_RUN_AS_ROOT_MSG)
			os.Exit(1)
		}

		if err := checks.CheckBPFDependencies(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Check if lock exists (profiler already running)
		if helper.IsFileExist(paths.LOCKFILEPATH) {
			fmt.Printf("ERROR: Another scx-adapt profile is already running. (%s)\n", paths.LOCKFILEPATH)
			os.Exit(1)
		}

		// Check if sched_ext is already active
		if checks.IsSchedExtActive() {
			fmt.Println("ERROR: sched_ext is already active")
			os.Exit(1)
		}

		// Create DATAFOLDER folder if not exist
		if err := helper.CreateDirIfNotExist(paths.DATAFOLDER); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// If profile exists in PROFILESFOLDER with that name, use it
		if helper.IsFileExist(path.Join(paths.PROFILESFOLDER, filepath)) {
			filepath = path.Join(paths.PROFILESFOLDER, filepath)
		}

		yamlData, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Printf("ERROR: Reading file '%s': %s\n", filepath, err)
			os.Exit(1)
		}

		// Convert YAML file bytes to 'Config' struct
		conf, err := helper.YamlToConfig(yamlData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create lock file
		if err := helper.CreateLock(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Interrupt handling
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		stop := make(chan bool, 1)                     // To send stop signal to Scheduler.Run() function
		errmsg := make(chan error, 1)                  // To receive error messages from Config.Run() and Scheduler.Run()
		schedChanged := make(chan helper.Scheduler, 1) // If current viable scheduler changes in Config.Run(), it sends the new 'Scheduler' instance,

		go conf.Run(schedChanged, errmsg)

		// Profile started message
		fmt.Printf("INFO: Profile started: '%s'\n", filepath)

	STOPERROR: // If a scheduler stop fails, execution continues here with an error in 'errsmg' channel. It gets consumed in  'case err := <-errmsg:'
		for {
			select {
			case err := <-errmsg:
				fmt.Println("ERROR:", err)

				if e := helper.RemoveLock(); e != nil {
					fmt.Println(e)
				}

				os.Exit(1)

			// When the new scheduler gets consumed with '<-schedChanged', stop signal is sent to the previously running scheduler (if there is)
			case sched := <-schedChanged:
				switch sched.Path {
				case "":
					fmt.Println("INFO: None of sched_ext schedulers match criterias. Switching to system scheduler...")

				default:
					fmt.Printf("INFO: Criterias match for scheduler '%s'...\n", sched.Path)
				}

				if checks.IsSchedExtActive() {
					stop <- true
					if err := <-errmsg; err != nil {
						errmsg <- err
						goto STOPERROR
					}
				}

				// If the new scheduler is a null-scheduler (system scheduler, none of entries in configuration are viable) 'Scheduler.Run()' is skipped
				if sched.Path != "" {
					go sched.Run(stop, errmsg) // // Then the new scheduler runs with 'Scheduler.Run()'

					fmt.Printf("INFO: Starting scheduler '%s'...\n", sched.Path)
				}

			case <-interrupt:
				fmt.Println(INTERRUPT_MSG)

				if checks.IsSchedExtActive() {
					stop <- true
					if err := <-errmsg; err != nil {
						errmsg <- err
						goto STOPERROR
					}
				}

				if e := helper.RemoveLock(); e != nil {
					fmt.Println(e)
				}

				os.Exit(0)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startProfileCmd)
}
