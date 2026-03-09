package helper

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/dogukanmeral/scx-adapt/internal/checks"
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

// Removes files in '/sys/fs/bpf/sched_ext' if exists (stops currently running sched_ext scheduler).
func (s Scheduler) Stop() error {
	switch s.Type {
	case "kernel":
		err := os.RemoveAll("/sys/fs/bpf/sched_ext/")
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("No custom schedulers are attached")
		} else if err != nil {
			return fmt.Errorf("Error occured while stopping current scheduler: %s\n", err)
		}
	}

	return nil
}

// Attaches sched_ext scheduler to kernel using 'bpftool' at '/sys/fs/bpf/sched_ext'
func (s Scheduler) Start() error {
	switch s.Type {
	case "kernel":
		if err := checks.CheckDependencies(); err != nil {
			return err
		}

		if err := checks.CheckObj(s.Path); err != nil {
			return err
		}

		startCmd := exec.Command("bpftool", "struct_ops", "register", s.GetAbsolutePath(), "/sys/fs/bpf/sched_ext")
		err := startCmd.Run()

		if err != nil {
			return fmt.Errorf("Error occured while attaching scheduler '%s': %s\n", s.Path, err)
		}
	}

	return nil
}
