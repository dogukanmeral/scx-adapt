/*
Copyright © 2026 Doğukan Meral <dogukan.meral@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dogukanmeral/scx-adapt/internal/checks"
	"github.com/spf13/cobra"
)

// checkDependenciesCmd represents the checkDepencies command
var checkDependenciesCmd = &cobra.Command{
	Use:   "check-dependencies",
	Short: "Check dependencies of scx-adapt.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Too many arguments. scx-adapt --help to see usage")
			os.Exit(1)
		}

		var depNotFound bool = false

		if err := checks.IsBpfToolInstalled(); err != nil {
			fmt.Print(err)
			depNotFound = true
		}

		if err := checks.IsBpfFsMounted(); err != nil {
			fmt.Print(err)
			depNotFound = true
		}

		if err := checks.IsSchedExtDirExist(); err != nil {
			fmt.Print(err)
			depNotFound = true
		}

		if depNotFound {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkDependenciesCmd)
}
