package helper

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// To extract integer values from files which has the format: <variablename>: <value>
// Use for:
//
//	/proc/stat: procs_running, procs_blocked
func GetVariableAsInt(filePath string, variableName string) (int, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return -1, fmt.Errorf("Error occured while reading file '%s': %s\n", filePath, err)
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s", variableName))

	var v string
out:
	for _, line := range strings.Split(string(data), "\n") {
		if re.MatchString(line) {
			v = strings.Fields(line)[1]
			break out
		}
	}

	val, err := strconv.Atoi(v)

	if err != nil {
		return -1, fmt.Errorf("Error occured while converting '%s' to int: %s\n", v, err)
	}

	return val, nil
}

// Read system file '/proc/diskstats' and return number of process with current IO operations.
func DiskCurIO() (int, error) {
	diskData, err := os.ReadFile("/proc/diskstats")

	if err != nil {
		return -1, fmt.Errorf("Error occured while reading file '%s': %s\n", "/proc/diskstats", err)
	}

	var curIO int

	for line := range strings.Lines(string(diskData)) {
		if partitionHier := strings.Fields(line)[1]; partitionHier == "0" {

			curIOPartition, err := strconv.Atoi(strings.Fields(line)[11])

			if err != nil {
				return -1, fmt.Errorf("Error occured while converting '%s' to int: %s\n", strings.Fields(line)[11], err)
			}

			curIO += curIOPartition
		}
	}

	return curIO, nil
}

type LoadAvgMinute int

const (
	Avg1min  LoadAvgMinute = 0
	Avg5min  LoadAvgMinute = 1
	Avg15min LoadAvgMinute = 2
)

// Read system file '/proc/loadavg' and return the load average for specified minute value.
func LoadAvg(minutes LoadAvgMinute) (float64, error) {
	loadAvgData, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, fmt.Errorf("Error occured while reading file '%s': %s\n", "/proc/loadavg", err)
	}

	loadAvgsStr := strings.Fields(string(loadAvgData))[:3]

	loadAvgsFloat := make([]float64, len(loadAvgsStr))

	for i, s := range loadAvgsStr {
		loadAvgsFloat[i], err = strconv.ParseFloat(s, 64)

		if err != nil {
			return 0, fmt.Errorf("Error occured while converting '%s' to float64: %s\n", s, err)
		}
	}

	return loadAvgsFloat[minutes], nil
}

// Parse load average value name and return it's minute.
func ParseLoadAvg(loadAvgValName string) LoadAvgMinute {
	var laMinute LoadAvgMinute

	switch strings.Split(loadAvgValName, "_")[2] {
	case "1":
		laMinute = Avg1min
	case "5":
		laMinute = Avg5min
	case "15":
		laMinute = Avg15min
	}

	return laMinute
}

type PressureType string

const (
	Cpu PressureType = "cpu"
	IO  PressureType = "io"
	Mem PressureType = "memory"
)

type PressureOption string

const (
	Some PressureOption = "some"
	Full PressureOption = "full"
)

type PressureSecond int

const (
	Avg10sec  PressureSecond = 0
	Avg60sec  PressureSecond = 1
	Avg300sec PressureSecond = 2
)

/*
 	“some” line indicates the share of time in which at least some tasks are stalled on a given resource.
	“full” line indicates the share of time in which all non-idle tasks are stalled on a given resource simultaneously.
*/

// Read system file '/proc/pressure/<pressure type>' and return the pressure for specified type, option and second values.
func Pressure(presType PressureType, presOpt PressureOption, presSec PressureSecond) (float64, error) {
	presFile := fmt.Sprintf("/proc/pressure/%s", presType)

	presData, err := os.ReadFile(presFile)
	if err != nil {
		return 0, fmt.Errorf("Error occured while reading file '%s': %s\n", presFile, err)
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s", presOpt))

	var pressures []float64
	for line := range strings.Lines(string(presData)) {
		if re.MatchString(line) {
			for _, psi := range strings.Fields(line)[1:4] {
				presStr := strings.Split(psi, "=")[1]

				presFloat64, err := strconv.ParseFloat(presStr, 64)

				if err != nil {
					return 0, fmt.Errorf("Error occured while converting '%s' to float64: %s\n", presStr, err)
				}

				pressures = append(pressures, presFloat64)
			}
		}
	}

	return pressures[presSec], nil
}

// Parse pressure value name to it's: type (CPU, IO or Mem), option(some or full), seconds(10, 60 or 300).
func ParsePressure(pressureValName string) (PressureType, PressureOption, PressureSecond) { // 0: avg10, 1: avg60, 2: avg300
	var pType PressureType
	var pOpt PressureOption
	var pSec PressureSecond

	switch strings.Split(pressureValName, "_")[0] {
	case string(Cpu):
		pType = Cpu
	case string(IO):
		pType = IO
	case string(Mem):
		pType = Mem
	}

	switch strings.Split(pressureValName, "_")[2] {
	case string(Some):
		pOpt = Some
	case string(Full):
		pOpt = Full
	}

	switch strings.Split(pressureValName, "_")[3] {
	case "10":
		pSec = Avg10sec
	case "60":
		pSec = Avg60sec
	case "300":
		pSec = Avg300sec
	}

	return pType, pOpt, pSec
}

// Convert elements in []float64 to string and return these as []string
func FloatsToStr(slice []float64) []string {

	out := make([]string, len(slice))

	for i, v := range slice {
		out[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}

	return out
}
