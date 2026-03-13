package helper

import (
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

func (s Scheduler) Run(stop <-chan bool, errmsg chan<- error) {
	var cmd *exec.Cmd

	switch s.Type {
	case string(KernelOnly):
		cmd = exec.Command("bpftool", "struct_ops", "register", s.GetAbsolutePath(), "/sys/fs/bpf/sched_ext")

	case string(Userspace):
		if s.Parameters != nil {
			cmd = exec.Command(s.GetAbsolutePath(), *s.Parameters...)
		} else {
			cmd = exec.Command(s.GetAbsolutePath())
		}
	}

	if err := cmd.Start(); err != nil {
		errmsg <- err
		return
	}

	finished := make(chan error, 1)

	go func() {
		finished <- cmd.Wait()
	}()

SELECTSTART:
	select {
	case err := <-finished:
		if err != nil {
			errmsg <- err
		} else {
			goto SELECTSTART
		}

	case <-stop:
		switch s.Type {
		case string(KernelOnly):
			if err := os.RemoveAll("/sys/fs/bpf/sched_ext/"); err != nil {
				errmsg <- fmt.Errorf("Error occured while detaching kernel-only scheduler '%s': %s\n", s.GetAbsolutePath(), err)
			} else {
				errmsg <- nil
			}
		case string(Userspace):
			if err := cmd.Process.Kill(); err != nil {
				errmsg <- fmt.Errorf("Error occured while stopping userspace scheduler '%s': %s\n", s.GetAbsolutePath(), err)
			} else {
				errmsg <- nil
			}
		}
	}
}
