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

type SchedulerLoader string

const (
	External SchedulerLoader = "external"
	Builtin  SchedulerLoader = "builtin"
)

type Scheduler struct {
	Loader     string     `yaml:"loader" validate:"required"`
	Parameters *[]string  `yaml:"parameters"`
	Log        *bool      `yaml:"log"`
	StructName *string    `yaml:"structName"`
	Path       string     `yaml:"path" validate:"required"`
	Priority   int        `yaml:"priority" validate:"required,gte=1,lte=139"`
	Criterias  []Criteria `yaml:"criterias" validate:"required,dive"`
}

var NilScheduler Scheduler = Scheduler{"", nil, nil, nil, "", 0, nil}

// Returns: path as it is if an absolute path, if not path of scheduler in SCHEDULERSFOLDER if exists, if none of both path as it is
func (s Scheduler) GetAbsolutePath() string {
	if path.IsAbs(s.Path) {
		return s.Path
	}

	var subdir string

	switch s.Loader {
	case string(External):
		subdir = paths.EXTERNALFOLDER
	case string(Builtin):
		subdir = paths.BUILTINFOLDER
	}

	if p := path.Join(subdir, s.Path); IsFileExist(p) {
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

	// Check if scheduler loader is valid (external, builtin)
	if !slices.Contains([]string{string(External), string(Builtin)}, s.Loader) {
		return &errs.InvalidSchedulerLoaderError{
			Msg: fmt.Sprintf("Invalid scheduler loader '%s' for scheduler '%s'.", s.Loader, s.Path),
		}
	}

	// Check if log option used only for builtin-loader schedulers
	if s.Log != nil {
		if s.Loader != string(Builtin) && *s.Log == true {
			return &errs.LogForNonBuiltinLoaderSchedError{
				Msg: fmt.Sprintf("Logging feature is not available for loader type '%s' for scheduler '%s'", s.Loader, s.Path),
			}
		}
	}

	// Check if structName is set
	if s.StructName != nil {
		if s.Loader != string(External) {
			return &errs.StructNameForNonExternalLoaderSchedError{
				Msg: fmt.Sprintf("Scheduler '%s' is not externally-loaded but 'structName' is set", s.Path),
			}
		}
	} else if s.Loader == string(External) {
		return &errs.StructNameNotSetError{
			Msg: fmt.Sprintf("Struct name is not set for externally-loaded scheduler %s", s.Path),
		}
	}

	// Check if parameters section is valid
	if s.Loader == string(External) && s.Parameters != nil {
		return &errs.ParametersForExternalSchedError{
			Msg: fmt.Sprintf("Runtime parameters cannot be passed to externally-loaded scheduler '%s'", s.Path),
		}
	}

	// If scheduler loader is external, check if file at the path exists and a BPF object file
	// If scheduler loader is builtin, chech if file at the path exists and is an executable file
	switch s.Loader {
	case string(External):
		if err := checks.CheckObj(s.GetAbsolutePath()); err != nil {
			return err
		}
	case string(Builtin):
		if !IsFileExist(s.GetAbsolutePath()) {
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
	errOut := make(chan error, 1) // To receive error messages from subroutine, where stdout and stderr is read

	var logActive bool
	if s.Log != nil { // Since log parameter is optional, this statement prevents null pointer deference
		logActive = *s.Log
	} else {
		logActive = false
	}

	switch s.Loader {
	case string(External):
		if err := LoadBPFScx(s.GetAbsolutePath(), *s.StructName); err != nil { // Use scx-adapt's BPF loader if loader type is external
			errOut <- err
		}

	case string(Builtin):
		if s.Parameters != nil { // Since command parameters are optional, this statement prevents null pointer deference
			cmd = exec.Command(s.GetAbsolutePath(), *s.Parameters...)
		} else {
			cmd = exec.Command(s.GetAbsolutePath())
		}

		if logActive {
			logFile, err := CreateLogFile(path.Base(s.Path))
			if err != nil {
				fmt.Printf("WARNING: %s\n", err)
				goto skipLogging // Do not stop execution of scheduler if log file creation fails
			}

			defer logFile.Close()

			cmd.Stdout = logFile
			cmd.Stderr = logFile

		skipLogging:

			if err := cmd.Start(); err != nil {
				errmsg <- err
				return
			}

			go func() {
				errOut <- cmd.Wait() // If a builtin-loader (executable) fails, error is sent to errOut channel, received in for loop
			}()
		}
	}

	for {
		select {
		case err := <-errOut:
			switch s.Loader {
			case string(Builtin):
				if logActive {
					errmsg <- fmt.Errorf("Scheduler '%s' exited unexpectedly, logs are available at: %s", s.Path, paths.LOGFOLDER)
				} else {
					errmsg <- fmt.Errorf("Scheduler '%s' exited unexpectedly", s.Path)
				}
				return

			case string(External):
				if err != nil {
					errmsg <- err
					return
				}

				// If BPF is loaded and error message is nil, continue
				// TODO: Add: Checks for linked scheduler's health
			}

		case <-stop:
			switch s.Loader {
			case string(External):
				if err := os.Remove(paths.SCHEDBPFPINPATH); err != nil {
					errmsg <- fmt.Errorf("Detaching scheduler '%s': %s\n", s.GetAbsolutePath(), err)
				} else {
					errmsg <- nil
				}

				return

			case string(Builtin):
				if err := cmd.Process.Kill(); err != nil {
					errmsg <- fmt.Errorf("Stopping scheduler '%s': %s\n", s.GetAbsolutePath(), err)
				} else {
					errmsg <- nil
				}

				return
			}
		}
	}
}
