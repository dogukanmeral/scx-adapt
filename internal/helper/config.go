package helper

import (
	"bytes"
	"fmt"

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
