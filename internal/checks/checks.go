package checks

import (
	"debug/elf"
	"errors"
	"fmt"
	"internal/errs"
	"os"
	"os/exec"
	"slices"
)

// Checks if file is BPF bytecode and includes .struct_ops.link section
func CheckObj(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return err // check for os.ErrNotExist in tests
	}

	elfFile, err := elf.NewFile(file)
	if err != nil {
		return &errs.NotObjFileError{Msg: fmt.Sprintf("Not an object file: %s\n", path)}
	}

	if elfFile.Type != elf.ET_REL || elfFile.Machine != elf.EM_BPF {
		return &errs.NotBPFFileError{Msg: fmt.Sprintf("Not a BPF file: %s\n", path)}
	}

	hasStructOpsLink := false

	for _, sec := range elfFile.Sections {
		switch sec.Name {
		case ".struct_ops.link":
			hasStructOpsLink = true
		}
	}

	if !hasStructOpsLink {
		return &errs.NoStructOpsError{Msg: fmt.Sprintf("Doesn't include '.struct_ops.link' section: %s\n", path)}
	}

	return nil
}

// Checks dependencies: bpftool, kernel (BPF and sched_ext)
func CheckDependencies() error {
	// Check if BPF tool is installed
	whichCmd := exec.Command("which", "bpftool")
	whichCmd.Run()
	err := whichCmd.Err

	if err != nil {
		return fmt.Errorf("'bpftool' is not found in PATH: %s\n", os.Getenv("PATH"))
	}

	// Check if BPF filesystem is mounted
	grepBpfCmd := "mount | grep bpf"
	err = exec.Command("bash", "-c", grepBpfCmd).Err

	if err != nil {
		return fmt.Errorf("Kernel does not have BPF functionalities")
	}

	// Check if sched_ext directory exists
	_, err = os.Stat("/sys/kernel/sched_ext")
	if os.IsNotExist(err) {
		return fmt.Errorf("Kernel does not have sched_ext functionalities")
	}

	return nil
}

// Checks if "/sys/kernel/sched_ext/root/ops" exists (if sched_ext is active).
func IsScxRunning() bool {
	opsFile := "/sys/kernel/sched_ext/root/ops"
	_, err := os.Stat(opsFile)

	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// If passed slice of comparables has duplicates inside:
//
//	returns true, slice of duplicating values
//
// else:
//
//	returns false, nil
func ContainsDuplicate[T comparable](arr []T) (bool, []T) {
	var duplicates []T

	for i, a := range arr {
		if slices.Contains(arr[i+1:], a) {
			duplicates = append(duplicates, a)
		}
	}

	if len(duplicates) == 0 {
		return false, nil
	} else {
		return true, duplicates
	}
}

// Returns if file exists or not
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
