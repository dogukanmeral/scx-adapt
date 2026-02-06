/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"internal/helper"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// BUG: Fix problems arised by changes in the "telemetry.go"

// logCsvCmd represents the log-csv command
var logCsvCmd = &cobra.Command{
	Use:   "log-csv",
	Short: "Print system variables to file in csv format",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var filepath string
		var interval float64

		switch len(args) {
		case 0:
			fmt.Println("Missing arguments. scx-adapt --help to see usage")
			os.Exit(1)
		case 1:
			filepath = args[0]
			interval = 1 // second
		case 2:
			filepath = args[0]
			if i, err := strconv.ParseFloat(args[1], 64); err != nil {
				fmt.Println("Error: Interval argument must be a positive integer.")
				os.Exit(1)
			} else {
				interval = i
			}
		default:
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}

		f, err := os.Create(filepath)

		if err != nil {
			fmt.Printf("Error occured while creating file '%s': %s\n", filepath, err)
			os.Exit(1)
		}

		features := []string{
			"time_ms",
			"cpu_psi_some_10",
			"cpu_psi_some_60",
			"cpu_psi_some_300",
			"cpu_psi_full_10",
			"cpu_psi_full_60",
			"cpu_psi_full_300",
			"io_psi_some_10",
			"io_psi_some_60",
			"io_psi_some_300",
			"io_psi_full_10",
			"io_psi_full_60",
			"io_psi_full_300",
			"mem_psi_some_10",
			"mem_psi_some_60",
			"mem_psi_some_300",
			"mem_psi_full_10",
			"mem_psi_full_60",
			"mem_psi_full_300",
			"load_avg_1",
			"load_avg_5",
			"load_avg_15",
			"procs_running",
			"procs_blocked",
			"procs_disk_io",
		}

		// First line (column names)
		featuresLine := strings.Join(features, ",")
		_, err = f.WriteString(fmt.Sprintf("%s\n", featuresLine))

		if err != nil {
			fmt.Printf("Error occured while writing features line to file '%s': %s", filepath, err)
			os.Exit(1)
		}

		// Interrupt handling (CTRL + Z)
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-ch
			fmt.Println("Exiting...")
			f.Close()
			os.Exit(0)
		}()

		prTypes := []helper.PressureType{
			helper.Cpu,
			helper.IO,
			helper.Mem,
		}

		prOpts := []helper.PressureOption{
			helper.Some,
			helper.Full,
		}

		prSeconds := []helper.PressureSecond{
			helper.Avg10sec,
			helper.Avg60sec,
			helper.Avg300sec,
		}

		laMinutes := []helper.LoadAvgMinute{
			helper.Avg1min,
			helper.Avg5min,
			helper.Avg15min,
		}

		buf := make([]string, 0, len(features))

		var curTime float64 = 0

		for {
			// Current time after start (milliseconds)
			buf = append(buf, strconv.FormatFloat(curTime, 'f', -1, 64))

			// Iterate over all pressures
			for _, t := range prTypes {
				for _, o := range prOpts {
					for _, s := range prSeconds {
						v, err := helper.Pressure(t, o, s)
						if err != nil {
							fmt.Println("Error occured while reading pressures.")
							os.Exit(1)
						}

						buf = append(buf, strconv.FormatFloat(v, 'f', -1, 64))
					}
				}
			}

			// Load averages
			for _, m := range laMinutes {
				v, err := helper.LoadAvg(m)

				if err != nil {
					fmt.Println("Error occured while reading load averages.")
					os.Exit(1)
				}

				buf = append(buf, strconv.FormatFloat(v, 'f', -1, 64))
			}

			// Processes
			if procsR, err := helper.GetVariableAsInt("/proc/stat", "procs_running"); err != nil {
				fmt.Println("Error occured while reading procs_running.")
				os.Exit(1)
			} else {
				buf = append(buf, strconv.Itoa(procsR))
			}

			if procsB, err := helper.GetVariableAsInt("/proc/stat", "procs_blocked"); err != nil {
				fmt.Println("Error occured while reading procs_blocked.")
				os.Exit(1)
			} else {
				buf = append(buf, strconv.Itoa(procsB))
			}

			if procsIO, err := helper.DiskCurIO(); err != nil {
				fmt.Println("Error occured while reading diskstats.")
				os.Exit(1)
			} else {
				buf = append(buf, strconv.Itoa(procsIO))
			}

			// Write row
			if len(buf) != 0 {
				row := strings.Join(buf, ",")
				_, err := f.WriteString(row + "\n")

				if err != nil {
					fmt.Printf("Error occured while writing to file '%s': %s\n", filepath, err)
					os.Exit(1)
				}
			}

			buf = []string{}

			time.Sleep(time.Second * time.Duration(interval)) // BUG: Below zero values e.g. 0.5
			curTime += interval * 1000
		}
	},
}

func init() {
	rootCmd.AddCommand(logCsvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCsvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCsvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
