package helper

import (
	"errors"
	"fmt"
	"internal/checks"
	"os"
	"os/exec"
)

// NOTE: No helper depends on another (except Write()), combine them in cmd and config logic

func CurrentScx() error {
	opsFile := "/sys/kernel/sched_ext/root/ops"

	if _, err := os.Stat(opsFile); err == nil {
		data, err := os.ReadFile(opsFile)

		if err != nil {
			return fmt.Errorf("Error occured while reading '%s'.\n", opsFile)
		}

		fmt.Printf("%s", string(data))

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No custom schedulers are attached")
	}

	return nil
}

func Write(path string, data string) {
	err := os.WriteFile(path, []byte(data), 0644)

	if err != nil {
		panic(err)
	}
}

func TraceSchedExt() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("Must run as root")
	}

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

	buf := make([]byte, 4096) // heap allocation

	for {
		n, err := f.Read(buf)
		if err != nil {
			return err
		}
		fmt.Print(string(buf[:n]))
	}
}

func StopCurrScx() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("Must run as root")
	}

	err := os.Remove("/sys/fs/bpf/sched_ext/sched_ops")
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("No custom schedulers are attached")
	} else if err != nil {
		return fmt.Errorf("Error occured while stopping current scheduler: %s\n", err)
	}

	return nil
}

func StartScx(scxPath string) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("Must run as root")
	}

	if err := checks.CheckDependencies(); err != nil {
		return err
	}

	if err := checks.CheckObj(scxPath); err != nil {
		return err
	}

	startCmd := exec.Command("bpftool", "struct_ops", "register", scxPath, "/sys/fs/bpf/sched_ext")
	startCmd.Run()

	if err := startCmd.Err; err != nil {
		return fmt.Errorf("Error occured while attaching scheduler: %s\n", err)
	}

	return nil
}
