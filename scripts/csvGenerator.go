package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getLoadAvgs() []string {
	loadavgData, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		panic(err)
	}

	return strings.Fields(string(loadavgData))[:3]
}

func sumOfArr(arr []int) int {
	sum := 0
	for _, val := range arr {
		sum += val
	}
	return sum
}

func getTotalIOWait() string {
	statData, err := os.ReadFile("/proc/stat")
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("^cpu")

	var iowaits []int

	for _, line := range strings.Split(string(statData), "\n") {
		if re.MatchString(line) {
			iowait, err := strconv.Atoi(strings.Fields(line)[5])
			if err != nil {
				panic(err)
			}

			iowaits = append(iowaits, iowait)
		}
	}

	return strconv.Itoa(sumOfArr(iowaits))
}

// for files like: variablename: value
// use for /proc/stats: ctxt, procs_running, procs_blocked
func getVal(filePath string, valName string) string {
	statData, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s", valName))

	var v string
out:
	for _, line := range strings.Split(string(statData), "\n") {
		if re.MatchString(line) {
			v = strings.Fields(line)[1]
			break out
		}
	}

	return v
}

type PressureType string

const (
	Cpu    PressureType = "cpu"
	IO     PressureType = "io"
	Memory PressureType = "memory"
)

type PressureOption string

const (
	Some PressureOption = "some"
	Full PressureOption = "full"
)

func getPressures(psiType PressureType, psiOpt PressureOption) []string {
	pressureFile := fmt.Sprintf("/proc/pressure/%s", psiType)

	pressureData, err := os.ReadFile(pressureFile)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s", psiOpt))

	var psis []string
	for line := range strings.Lines(string(pressureData)) {
		if re.MatchString(line) {
			for _, psi := range strings.Fields(line)[1:4] {
				psis = append(psis, strings.Split(psi, "=")[1])
			}
		}
	}

	return psis
}

func getDiskstats() []string {
	diskstatsData, err := os.ReadFile("/proc/diskstats")
	if err != nil {
		panic(err)
	}

	var partitionsIOcur int
	var partitionsIOms int

	for line := range strings.Lines(string(diskstatsData)) {
		if partitionHier := strings.Fields(line)[1]; partitionHier == "0" {

			curIO, err := strconv.Atoi(strings.Fields(line)[11])

			if err != nil {
				panic(err)
			}

			partitionsIOcur += curIO

			msIO, err := strconv.Atoi(strings.Fields(line)[12])

			if err != nil {
				panic(err)
			}

			partitionsIOms += msIO
		}
	}

	return []string{strconv.Itoa(partitionsIOcur), strconv.Itoa(partitionsIOms)}
}

func getNetDev() []string {
	netDevData, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		panic(err)
	}

	var rtxPackets int
	var txPackets int

	var curLine int = 0

	for line := range strings.Lines(string(netDevData)) {
		if curLine < 2 {
			curLine += 1
			continue
		}

		rtx, err := strconv.Atoi(strings.Fields(line)[2])
		if err != nil {
			panic(err)
		}
		rtxPackets += rtx

		tx, err := strconv.Atoi(strings.Fields(line)[10])
		if err != nil {
			panic(err)
		}
		txPackets += tx
	}

	return []string{strconv.Itoa(rtxPackets), strconv.Itoa(txPackets)}
}

var features []string = []string{"loadavg-1min", "loadavg-5min", "loadavg-15min", "iowait-ms",
	"ctxt", "procs-r", "procs-b", "diskcurIO", "diskIOms", "rtxPackets", "txPackets"}

func main() {
	csvFile, err := os.Create("sys-data.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write(features)

	for {
		loadAvgs := getLoadAvgs()
		totalIOwait := getTotalIOWait()
		diskStats := getDiskstats()
		netDev := getNetDev()

		ctxt := getVal("/proc/stat", "ctxt")
		procsR := getVal("/proc/stat", "procs_running")
		procsB := getVal("/proc/stat", "procs_blocked")

		writer.Write([]string{loadAvgs[0], loadAvgs[1], loadAvgs[2], totalIOwait,
			ctxt, procsR, procsB, diskStats[0], diskStats[1], netDev[0], netDev[1]})

		writer.Flush()
		time.Sleep(time.Second)
	}
}
