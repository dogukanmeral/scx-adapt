package helper

import (
	"fmt"
	"regexp"

	"github.com/dogukanmeral/scx-adapt/internal/errs"
	"github.com/go-playground/validator/v10"
)

type Criteria struct {
	ValueName string   `yaml:"value_name" validate:"required"`
	MoreThan  *float64 `yaml:"more_than"`
	LessThan  *float64 `yaml:"less_than"`
}

/*
	Valid value_name(s):
		(cpu|io|mem)_psi_(some|full)_(10|60|300)
		load_avg_(1|5|15)
		procs_running
		procs_blocked
		procs_disk_io
*/

var VALID_VALUE_REGEX = map[string]string{
	"pressures":    "^(cpu|io|mem)_psi_(some|full)_(10|60|300)$",
	"loadAvgs":     "^load_avg_(1|5|15)$",
	"procsRunning": "^procs_running$",
	"procsBlocked": "^procs_blocked$",
	"procsDiskIo":  "^procs_disk_io$",
}

// Validate Criteria
func (c Criteria) Validate() error {
	v := validator.New()

	if err := v.Struct(c); err != nil {
		return err
	}

	for _, r := range VALID_VALUE_REGEX {
		if m, _ := regexp.MatchString(r, c.ValueName); m {
			goto valueNameValid
		}
	}
	return &errs.InvalidValueNameError{
		Msg: fmt.Sprintf("Invalid value_name: %s", c.ValueName),
	}

valueNameValid:

	if c.MoreThan == nil && c.LessThan == nil {
		return &errs.MissingParameterError{
			Msg: fmt.Sprintf("There is no 'more_than' and/or 'less_than' parameter for value '%s'", c.ValueName),
		}
	}

	if c.MoreThan != nil && c.LessThan != nil {
		if *c.MoreThan >= *c.LessThan {
			return &errs.ConflictParametersError{
				Msg: fmt.Sprintf("Parameter 'more_than' cannot be >= 'less_than' in value '%s'", c.ValueName),
			}
		}
	}

	return nil
}

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
