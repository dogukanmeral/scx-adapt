package helper

import "regexp"

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
