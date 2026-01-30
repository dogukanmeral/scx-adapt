package helper

import (
	"errors"
	"fmt"
	"internal/checks"
	"os"
	"os/exec"
	"path/filepath"
)

// NOTE: No helper depends on another (except Write()), combine them in cmd and config logic

func AddScx(scxPath string, addedScxsPath string) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("Must run as root")
	}

	// Check object
	if err := checks.CheckObj(scxPath); err != nil {
		return err
	}

	// Check if same file exists
	curScheds, err := os.ReadDir(addedScxsPath)

	if err != nil {
		return fmt.Errorf("Error reading directory '%s': %s\n", addedScxsPath, err)
	}

	for _, e := range curScheds {
		if e.Name() == filepath.Base(scxPath) {
			return fmt.Errorf("Scheduler exists with same name: %s\n", filepath.Base(scxPath))
		}
	}

	// Copy file
	input, err := os.ReadFile(scxPath)
	if err != nil {
		return fmt.Errorf("Error occured while reading '%s': %s\n", scxPath, err)
	}

	copyFilePath := filepath.Join(addedScxsPath, filepath.Base(scxPath))
	err = os.WriteFile(copyFilePath, input, 0744)

	if err != nil {
		return fmt.Errorf("Error occured while writing '%s': %s\n", copyFilePath, err)
	}

	return nil
}

func RemoveAddedScx(scxFilename string, addedScxsPath string) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("Must run as root")
	}

	scheds, err := os.ReadDir(addedScxsPath)
	if err != nil {
		return fmt.Errorf("Error reading directory '%s': %s\n", addedScxsPath, err)
	}

	for _, s := range scheds {
		if s.Name() == scxFilename {
			err := os.Remove(filepath.Join(addedScxsPath, scxFilename))
			if err != nil {
				return fmt.Errorf("Error occured while removing scheduler '%s': %s\n", scxFilename, err)
			}

			return nil
		}
	}

	return fmt.Errorf("Scheduler does not exist: %s\n", scxFilename)
}

func ListScxs(addedScxsPath string) error {
	files, err := os.ReadDir(addedScxsPath)

	if err != nil {
		return fmt.Errorf("Error reading directory '%s': %s\n", addedScxsPath, err)
	}

	for _, e := range files {
		fmt.Println(e.Name())
	}

	return nil
}

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

func StartScx(scxPath string, addedScxsPath string) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("Must run as root")
	}

	err := checks.CheckDependencies()
	if err != nil {
		return err
	}

	if filepath.IsAbs(scxPath) { // BUG: doesn't work on ../ ./ etc.
		err := checks.CheckObj(scxPath)
		if err != nil {
			return nil
		}

	} else {
		err := checks.CheckScxAdded(scxPath, addedScxsPath)
		if err != nil {
			return err
		}

		scxPath = filepath.Join(addedScxsPath, scxPath)
	}

	startCmd := exec.Command("bpftool", "struct_ops", "register", scxPath, "/sys/fs/bpf/sched_ext")
	startCmd.Run()
	err = startCmd.Err

	if err != nil {
		return fmt.Errorf("Error occured while attaching scheduler: %s\n", err)
	}

	return nil
}
