package helper

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/dogukanmeral/scx-adapt/internal/checks"
)

// Parses the profile and iterates on schedulers until all criterias of a scheduler is satisfied; then sleeps for a specified time period
func RunProfile(profilePath string) error {
	profileData, err := os.ReadFile(profilePath)
	if err != nil {
		return fmt.Errorf("Error occured while reading file '%s': %s\n", profilePath, err)
	}

	conf, err := YamlToConfig(profileData)
	if err != nil {
		return err
	}

	fmt.Printf("Profile at '%s' started.\n", profilePath)

	sort.Sort(conf) // Sort schedulers by their priority (smaller int has higher priority)

	var currentSched Scheduler

NEXT_SCHED:
	for i, s := range conf.Schedulers {
		for _, c := range s.Criterias {
			if b, err := c.Satisfies(); !b {
				if err != nil {
					return err
				}

				if i+1 == len(conf.Schedulers) {
					if checks.IsScxRunning() {
						err := StopCurrScx()

						if err != nil {
							return err
						}

						fmt.Println("None of schedulers satisfy criterias. Switched to system scheduler.")
						currentSched = Scheduler{"", 0, []Criteria{}}
					}
				}
				continue NEXT_SCHED
			}
		}

		if currentSched.Path != "" && !checks.IsSchedExtActive() {
			return fmt.Errorf("Error: Scheduler '%s' crashed.", currentSched.Path)
		}

		if s.Path != currentSched.Path {
			if checks.IsScxRunning() {
				err := StopCurrScx()

				if err != nil {
					return err
				}
			}

			err := StartScx(s.GetAbsolutePath())
			if err != nil {
				return err
			}

			fmt.Printf("Switched to scheduler '%s'\n", s.GetAbsolutePath())
			currentSched = s
		}

		goto SCHED_STARTED
	}

SCHED_STARTED:
	time.Sleep(time.Millisecond * time.Duration(conf.Interval))
	goto NEXT_SCHED
}
