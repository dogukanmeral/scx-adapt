package helper

import (
	"fmt"
	"internal/checks"
	"os"
	"regexp"
	"sort"
	"time"
)

// Checks if the system value satisfies 'less_than' or 'more_than'
func (c Criteria) SatisfiesLessMore(sysValue float64) bool {
	// Checking pointers to avoid null-pointer referance
	if c.MoreThan != nil && c.LessThan != nil {
		return sysValue > *c.MoreThan && sysValue < *c.LessThan
	} else if c.MoreThan != nil && c.LessThan == nil {
		return sysValue > *c.MoreThan
	} else if c.LessThan != nil && c.MoreThan == nil {
		return sysValue < *c.LessThan
	} else {
		return false
	}
}

// Checks if live system values satisfies the criteria.
func (c Criteria) Satisfies() (bool, error) {
	if b, _ := regexp.MatchString(VALID_VALUE_REGEX["pressures"], c.ValueName); b {
		pType, pOpt, pSec := ParsePressure(c.ValueName)
		pValue, err := Pressure(pType, pOpt, pSec)

		if err != nil {
			return false, err
		}

		return c.SatisfiesLessMore(pValue), nil

	} else if b, _ := regexp.MatchString(VALID_VALUE_REGEX["loadAvgs"], c.ValueName); b {
		laMinute := ParseLoadAvg(c.ValueName)
		laValue, err := LoadAvg(laMinute)

		if err != nil {
			return false, err
		}

		return c.SatisfiesLessMore(laValue), nil

	} else if b, _ := regexp.MatchString(VALID_VALUE_REGEX["procsRunning"], c.ValueName); b {
		pRunValue, err := GetVariableAsInt("/proc/stat", "procs_running")

		if err != nil {
			return false, err
		}

		return c.SatisfiesLessMore(float64(pRunValue)), nil

	} else if b, _ := regexp.MatchString(VALID_VALUE_REGEX["procsBlocked"], c.ValueName); b {
		pBlckValue, err := GetVariableAsInt("/proc/stat", "procs_blocked")

		if err != nil {
			return false, err
		}

		return c.SatisfiesLessMore(float64(pBlckValue)), nil

	} else if b, _ := regexp.MatchString(VALID_VALUE_REGEX["procsDiskIo"], c.ValueName); b {
		pIoValue, err := DiskCurIO()

		if err != nil {
			return false, err
		}

		return c.SatisfiesLessMore(float64(pIoValue)), nil
	} else {
		return false, nil
	}
}

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

		if s.Path != currentSched.Path {
			if checks.IsScxRunning() {
				err := StopCurrScx()

				if err != nil {
					return err
				}
			}

			err := StartScx(s.Path)
			if err != nil {
				return err
			}

			fmt.Printf("Switched to scheduler '%s'\n", s.Path)
			currentSched = s
		}

		goto SCHED_STARTED
	}

SCHED_STARTED:
	time.Sleep(time.Millisecond * time.Duration(conf.Interval))
	goto NEXT_SCHED
}
