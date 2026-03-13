/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dogukanmeral/scx-adapt/internal/helper"

	"github.com/spf13/cobra"
)

var features = []string{
	"time_ms",
	"cpu_cores",
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

var prTypes = []helper.PressureType{
	helper.Cpu,
	helper.IO,
	helper.Mem,
}

var prOpts = []helper.PressureOption{
	helper.Some,
	helper.Full,
}

var prSeconds = []helper.PressureSecond{
	helper.Avg10sec,
	helper.Avg60sec,
	helper.Avg300sec,
}

var laMinutes = []helper.LoadAvgMinute{
	helper.Avg1min,
	helper.Avg5min,
	helper.Avg15min,
}

// logCsvCmd represents the log-csv command
var logCsvCmd = &cobra.Command{
	Use:   "log-csv <csv_file_path> [interval_ms]",
	Short: "Print system variables to file in csv format",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var filepath string
		var interval time.Duration

		switch len(args) {
		case 0:
			fmt.Println(MISSING_ARGS_MSG)
			os.Exit(1)
		case 1:
			filepath = args[0]
			interval = 1000 // milliseconds
		case 2:
			filepath = args[0]
			if i, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("Error: Interval argument must be a positive integer.")
				os.Exit(1)
			} else {
				interval = time.Duration(i)
			}
		default:
			fmt.Println(TOO_MANY_ARGS_MSG)
			os.Exit(1)
		}

		f, err := os.Create(filepath)

		if err != nil {
			fmt.Printf("Error: Creating file '%s': %s\n", filepath, err)
			os.Exit(1)
		}

		// Interrupt and error handling (closes file)
		kill := make(chan os.Signal, 1)
		signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

		go func() {
			sig := <-kill

			f.Close()

			if slices.Contains([]os.Signal{os.Interrupt, syscall.SIGTERM}, sig) {
				fmt.Println(INTERRUPT_MSG)
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}()

		// First line (column names)
		featuresLine := strings.Join(features, ",")
		_, err = f.WriteString(fmt.Sprintf("%s\n", featuresLine))

		if err != nil {
			fmt.Printf("Error: Writing features line to file '%s': %s\n", filepath, err)
			kill <- os.Kill
		}

		buf := make([]string, 0, len(features))
		var curTime time.Duration = 0

		for {
			// Current time after start (milliseconds)
			buf = append(buf, strconv.Itoa(int(curTime)))

			// Total # of CPU cores
			c, err := helper.TotalCores()
			if err != nil {
				fmt.Println(err)
				kill <- os.Kill
			}

			buf = append(buf, strconv.Itoa(c))

			// Iterate over all pressures
			for _, t := range prTypes {
				for _, o := range prOpts {
					for _, s := range prSeconds {
						v, err := helper.Pressure(t, o, s)
						if err != nil {
							fmt.Printf("Error: Reading pressures: %s\n", err)
							kill <- os.Kill
						}

						buf = append(buf, strconv.FormatFloat(v, 'f', -1, 64))
					}
				}
			}

			// Load averages
			for _, m := range laMinutes {
				v, err := helper.LoadAvg(m)

				if err != nil {
					fmt.Printf("Error: Reading load averages: %s\n", err)
					kill <- os.Kill
				}

				buf = append(buf, strconv.FormatFloat(v, 'f', -1, 64))
			}

			// Processes
			if procsR, err := helper.GetVariableAsInt("/proc/stat", "procs_running"); err != nil {
				fmt.Printf("Error: Reading procs_running: %s\n", err)
				kill <- os.Kill
			} else {
				buf = append(buf, strconv.Itoa(procsR))
			}

			if procsB, err := helper.GetVariableAsInt("/proc/stat", "procs_blocked"); err != nil {
				fmt.Printf("Error: Reading procs_blocked: %s\n", err)
				kill <- os.Kill
			} else {
				buf = append(buf, strconv.Itoa(procsB))
			}

			if procsIO, err := helper.DiskCurIO(); err != nil {
				fmt.Printf("Error: Reading diskstats: %s\n", err)
				kill <- os.Kill
			} else {
				buf = append(buf, strconv.Itoa(procsIO))
			}

			// Write row
			if len(buf) != 0 {
				row := strings.Join(buf, ",")
				_, err := f.WriteString(row + "\n")

				if err != nil {
					fmt.Printf("Error: Writing to file '%s': %s\n", filepath, err)
					kill <- os.Kill
				}
			}

			buf = []string{}

			time.Sleep(time.Millisecond * interval)
			curTime += interval
		}
	},
}

func init() {
	rootCmd.AddCommand(logCsvCmd)
}
