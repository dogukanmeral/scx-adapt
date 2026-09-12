// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package profile

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dogukanmeral/scx-adapt/internal/helper"
	"github.com/dogukanmeral/scx-adapt/internal/msg"
	"github.com/dogukanmeral/scx-adapt/internal/paths"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit <profile-filename>",
		Short: "Edit a profile with the default editor",
		Long: `Edit a profile stored in the profiles folder.

Opens the profile in the default editor (from $VISUAL or $EDITOR). After the
editor closes, the profile is validated. If it is invalid, you are prompted to
restore the backup, edit again, or keep the changes.

Requires root privileges.`,
		Run: func(cmd *cobra.Command, args []string) {
			var profileName string

			switch len(args) {
			case 0:
				fmt.Println(msg.MISSING_ARGS_MSG)
				os.Exit(1)
			case 1:
				profileName = args[0]
			default:
				fmt.Println(msg.TOO_MANY_ARGS_MSG)
				os.Exit(1)
			}

			if os.Geteuid() != 0 {
				fmt.Println(msg.MUST_RUN_AS_ROOT_MSG)
				os.Exit(1)
			}

			profilePath := path.Join(paths.PROFILESFOLDER, profileName)

			// Check if profile exists in the profiles directory
			if !helper.IsFileExist(profilePath) {
				fmt.Printf("ERROR: Profile configuration with filename '%s' does not exist at '%s'\n",
					profileName, paths.PROFILESFOLDER)
				os.Exit(1)
			}

			// Keep a temporary backup in memory
			backup, err := os.ReadFile(profilePath)
			if err != nil {
				fmt.Printf("ERROR: Reading file '%s': %s\n", profilePath, err)
				os.Exit(1)
			}

			editor := resolveEditor()

			for {
				if err := runEditor(editor, profilePath); err != nil {
					fmt.Printf("ERROR: Editor '%s' failed: %s\n", editor, err)
					os.Exit(1)
				}

				// Validate the profile after the editor closes
				profileData, err := os.ReadFile(profilePath)
				if err != nil {
					fmt.Printf("ERROR: Reading file '%s': %s\n", profilePath, err)
					os.Exit(1)
				}

				if _, err = helper.YamlToConfig(profileData); err == nil {
					fmt.Printf("Valid config: %s\n", profilePath)
					return
				}

				fmt.Println(err)
				fmt.Println("The profile is invalid. What would you like to do?")
				fmt.Println("  [r] Restore the backup")
				fmt.Println("  [e] Edit again")
				fmt.Println("  [c] Continue (leave as is)")
				fmt.Print("> ")

				reader := bufio.NewReader(os.Stdin)
				choice, _ := reader.ReadString('\n')

				switch strings.ToLower(strings.TrimSpace(choice)) {
				case "r":
					if err := os.WriteFile(profilePath, backup, 0700); err != nil {
						fmt.Printf("ERROR: Restoring backup: %s\n", err)
						os.Exit(1)
					}
					fmt.Printf("Backup restored: %s\n", profilePath)
					return
				case "e":
					continue
				case "c":
					fmt.Printf("Leaving profile as is: %s\n", profilePath)
					return
				default:
					fmt.Println("Invalid choice. Leaving profile as is.")
					return
				}
			}
		},
	}
}

// resolveEditor returns the user's preferred editor from VISUAL/EDITOR env vars,
// falling back to "vi" if neither is set.
func resolveEditor() string {
	if e := os.Getenv("VISUAL"); e != "" {
		return e
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	return "vi"
}

// runEditor launches the editor on the given file, wiring the terminal through.
func runEditor(editor string, filePath string) error {
	execCmd := exec.Command(editor, filePath)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}
