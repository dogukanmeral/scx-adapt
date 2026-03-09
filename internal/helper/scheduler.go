package helper

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"slices"

	paths "github.com/dogukanmeral/scx-adapt/internal"
	"github.com/dogukanmeral/scx-adapt/internal/checks"
	"github.com/dogukanmeral/scx-adapt/internal/errs"
	"github.com/go-playground/validator/v10"
)

type SchedulerType string

const (
	KernelOnly SchedulerType = "kernelonly"
	Userspace  SchedulerType = "userspace"
)

type Scheduler struct {
	Type       string     `yaml:"type" validate:"required"`
	Parameters *[]string  `yaml:"parameters"`
	Path       string     `yaml:"path" validate:"required"`
	Priority   int        `yaml:"priority" validate:"required,gte=1,lte=139"`
	Criterias  []Criteria `yaml:"criterias" validate:"required,dive"`
}

// Returns: path as it is if an absolute path, if not path of scheduler in SCHEDULERSFOLDER if exists, if none of both path as it is
func (s Scheduler) GetAbsolutePath() string {
	if path.IsAbs(s.Path) {
		return s.Path
	} else if p := path.Join(paths.SCHEDULERSFOLDER, s.Path); checks.IsFileExist(p) {
		return p
	}

	return s.Path
}

// Validate Scheduler
func (s Scheduler) Validate() error {
	v := validator.New()

	if err := v.Struct(s); err != nil {
		return err
	}

	// Check if scheduler type is valid (kernelonly, userspace)
	if !slices.Contains([]string{string(KernelOnly), string(Userspace)}, s.Type) {
		return &errs.InvalidSchedulerTypeError{
			Msg: fmt.Sprintf("Invalid scheduler type '%s' for scheduler '%s'.", s.Type, s.Path),
		}
	}

	// Check if parameters section is valid
	if s.Type == string(KernelOnly) && s.Parameters != nil {
		return &errs.ParametersForKernelSchedError{
			Msg: fmt.Sprintf("Runtime parameters cannot be passed to kernel-only scheduler '%s'", s.Path),
		}
	}

	// If scheduler type is kernel-only, check if file at the path exists and a BPF object file
	// If scheduler type is userspace, chech if file at the path exists and is an executable file
	switch s.Type {
	case string(KernelOnly):
		if err := checks.CheckObj(s.GetAbsolutePath()); err != nil {
			return err
		}
	case string(Userspace):
		if !checks.IsFileExist(s.GetAbsolutePath()) {
			return &errs.SchedulerDoesNotExistError{
				Msg: fmt.Sprintf("Scheduler does not exist at path '%s'", s.GetAbsolutePath()),
			}
		} else if !checks.IsExecutableELF(s.GetAbsolutePath()) {
			return &errs.NotExecutableELFError{
				Msg: fmt.Sprintf("File at path '%s' is not an executable ELF", s.GetAbsolutePath()),
			}
		}
	}

	// Check all criterias inside scheduler
	var valueNames []string
	for _, c := range s.Criterias {
		valueNames = append(valueNames, c.ValueName)

		if err := c.Validate(); err != nil {
			return err
		}
	}

	// Check if a criteria is defined multiple times in same scheduler
	cont, dup := checks.ContainsDuplicate(valueNames)
	if cont {
		return &errs.ConflictCriteriasError{
			Msg: fmt.Sprintf("Criteria(s) '%s' defined multiple times for scheduler '%s'", dup, s.GetAbsolutePath()),
		}
	}

	return nil
}

// Removes files in '/sys/fs/bpf/sched_ext' if exists (stops currently running sched_ext scheduler).
func (s Scheduler) Stop() error {
	switch s.Type {
	case string(KernelOnly):
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
	case string(KernelOnly):
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
