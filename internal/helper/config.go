package helper

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/dogukanmeral/scx-adapt/internal/checks"
	"github.com/dogukanmeral/scx-adapt/internal/errs"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// Interface for sorting schedulers by their priority
func (c Config) Len() int {
	return len(c.Schedulers)
}

// Interface for sorting schedulers by their priority
func (c Config) Less(i, j int) bool {
	return c.Schedulers[i].Priority < c.Schedulers[j].Priority
}

// Interface for sorting schedulers by their priority
func (c Config) Swap(i, j int) {
	c.Schedulers[i], c.Schedulers[j] = c.Schedulers[j], c.Schedulers[i]
}

type Config struct {
	Interval   int         `yaml:"interval" validate:"required,gte=1"` // ms
	Schedulers []Scheduler `yaml:"schedulers" validate:"required,dive"`
}

// Validate Config
func (conf Config) Validate() error {
	v := validator.New()

	if err := v.Struct(conf); err != nil {
		return err
	}

	var priorities []int

	// Check all schedulers in config
	for _, s := range conf.Schedulers {
		priorities = append(priorities, s.Priority)

		if err := s.Validate(); err != nil {
			return err
		}
	}

	// Check if a priority is assigned to multiple schedulers
	cont, dup := checks.ContainsDuplicate(priorities)
	if cont {
		return &errs.ConflictPrioritiesError{Msg: fmt.Sprintf("Priority(s) '%d' is/are assigned for multiple schedulers", dup)}
	}

	return nil
}

// Converts YAML data passed as []byte to Config. If any error occurs in the called functions; returns it as it is.
func YamlToConfig(yamlData []byte) (Config, error) {
	var conf Config

	decoder := yaml.NewDecoder(bytes.NewReader(yamlData))
	decoder.KnownFields(true) // Check unrelated keys in YAML

	if err := decoder.Decode(&conf); err != nil {
		return conf, err
	}

	if err := conf.Validate(); err != nil {
		return conf, err
	}

	return conf, nil
}

func (conf Config) Run(changed chan<- Scheduler, errmsg chan<- error) {
	sort.Sort(conf) // Sort Config's by their priority (as defined in 'sort' interface of Config)

	var currentSched Scheduler = NilScheduler

NEXT_SCHED:
	for i, s := range conf.Schedulers {
		for _, c := range s.Criterias {
			if b, err := c.Satisfies(); !b { // If a criteria is not satisfied
				if err != nil { // If an error occurs while reading system variables, write the error to errmsg channel and return subroutine
					errmsg <- err
					return
				}

				if i+1 == len(conf.Schedulers) { // If it's the last scheduler in config
					if currentSched.Path != "" {
						currentSched = NilScheduler
						changed <- currentSched // Send nil-scheduler (system scheduler) to changed channel (gets received in start-profile)
					}
				}

				continue NEXT_SCHED // Continue to check criterias of next scheduler (end of for loop if last element)
			}
		}

		if s.Path != currentSched.Path { // Schedulers which satisfies all criterias reach here
			currentSched = s // Change current scheduler variable
			changed <- s     // Send new viable scheduler to changed channel (gets received in start-profile and new scheduler starts)
		}

		goto SCHED_STARTED // After switching to new viable scheduler, go to sleeping period
	}

SCHED_STARTED:
	time.Sleep(time.Millisecond * time.Duration(conf.Interval))
	goto NEXT_SCHED // At the end of sleeping period, go to start of for loop which iterates over schedulers and their criterias
}
