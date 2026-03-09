package helper

import (
	"errors"
	"fmt"
	"os"
)

// NOTE: No helper depends on another (except Write()), combine them in cmd and config logic

// Returns content of "/sys/kernel/sched_ext/root/ops" (currently running sched_ext scheduler if exists).
func CurrentScx() (string, error) {
	opsFile := "/sys/kernel/sched_ext/root/ops"

	if _, err := os.Stat(opsFile); err == nil {
		data, err := os.ReadFile(opsFile)

		if err != nil {
			return "", fmt.Errorf("Error occured while reading '%s'.\n", opsFile)
		}

		return string(data), nil

	} else if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("No custom schedulers are attached")
	} else {
		return "", err
	}
}

// Returns sched_ext trace (does not return anything if sched_ext is not active).
func TraceSchedExt(outfile string) error {
	// stop tracing
	Write("/sys/kernel/tracing/tracing_on", "0")

	// clear tracing data
	Write("/sys/kernel/tracing/trace", "")

	// enable sched_ext events
	Write("/sys/kernel/tracing/events/sched_ext/enable", "1")

	// start tracing
	Write("/sys/kernel/tracing/tracing_on", "1")

	defer Write("/sys/kernel/tracing/tracing_on", "0")
	defer Write("/sys/kernel/tracing/trace", "")

	f, err := os.Open("/sys/kernel/tracing/trace_pipe")
	if err != nil {
		return err
	}
	defer f.Close()

	o, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer o.Close()

	buf := make([]byte, 4096) // heap allocation

	for {
		nr, err := f.Read(buf)
		if err != nil {
			return err
		}

		_, err = o.Write(buf[:nr])
		if err != nil {
			return fmt.Errorf("Error occured while writing trace to file '%s': %s", outfile, err)
		}
	}
}
