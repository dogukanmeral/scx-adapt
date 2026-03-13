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

		if err := checks.CheckDependencies(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Check if lock exists (profiler already running)
		if checks.IsFileExist(paths.LOCKFILEPATH) {
			fmt.Printf("Error: Another scx-adapt profile is already running. (%s)\n", paths.LOCKFILEPATH)
			os.Exit(1)
		}

		// Create DATAFOLDER folder if not exist
		if err := helper.CreateDirIfNotExist(paths.DATAFOLDER); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// If profile exists in PROFILESFOLDER with that name, use it
		if checks.IsFileExist(path.Join(paths.PROFILESFOLDER, filepath)) {
			filepath = path.Join(paths.PROFILESFOLDER, filepath)
		}

		yamlData, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Printf("Error: Reading file '%s': %s\n", filepath, err)
			os.Exit(1)
		}

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

		stop := make(chan bool, 1)
		errmsg := make(chan error, 1)
		schedChanged := make(chan helper.Scheduler, 1)

		go conf.Run(schedChanged, errmsg)

	STOPERROR:
		for {
			select {
			case err := <-errmsg:
				fmt.Println(err)

				if e := helper.RemoveLock(); e != nil {
					fmt.Println(e)
				}

				os.Exit(1)

			case sched := <-schedChanged:
				switch sched.Path {
				case "":
					fmt.Println("None of sched_ext schedulers match criterias. Switching to system scheduler...")

				default:
					fmt.Printf("Switching to scheduler '%s'...\n", sched.Path)
				}

				if checks.IsSchedExtActive() {
					stop <- true
					if err := <-errmsg; err != nil {
						errmsg <- err
						goto STOPERROR
					}
				}

				if sched.Path != "" {
					go sched.Run(stop, errmsg)

					fmt.Printf("Starting scheduler '%s'...\n", sched.Path)
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
