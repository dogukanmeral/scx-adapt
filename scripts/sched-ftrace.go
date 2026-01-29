package main

import (
	"fmt"
	"os"
)

func write(path, val string) error {
	return os.WriteFile(path, []byte(val), 0644)
}

func main() {
	if os.Geteuid() != 0 { // root user id
		panic("must run as root")
	}

	base := "/sys/kernel/tracing"

	// stop previous tracing
	write(base+"/tracing_on", "0")

	// clear old tracing data
	write(base+"/trace", "")

	// enable all sched events
	write(base+"/events/sched/enable", "1")

	// start tracing
	write(base+"/tracing_on", "1")

	f, err := os.Open(base + "/trace_pipe")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := make([]byte, 4096) // heap allocation
	for {
		n, err := f.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Print(string(buf[:n]))
	}
}
