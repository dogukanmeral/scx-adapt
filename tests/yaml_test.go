package main

import (
	"errors"
	"fmt"
	"internal/errs"
	"internal/helper"
	"os"
	"testing"
)

func formatYamlError(yamlData string, err error) string {
	return fmt.Sprintf("\n\nYAML CONFIGURATION:\n\n%s\n\nUnexpected error: %v", yamlData, err)
}

////////// Positive tests //////////

func TestSingleSched(t *testing.T) {
	yamlData := `interval: 1000
schedulers:
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 5
        more_than: 0
      - value_name: procs_running
        more_than: 3
        less_than: 10`

	_, err := helper.YamlToConfig([]byte(yamlData))

	if err != nil {
		t.Error(formatYamlError(yamlData, err))
	}
}

func TestMultipleScheds(t *testing.T) {
	yamlData := `interval: 1000
schedulers:
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 5
  - path: "../../bytecode/lottery.bpf.c.o"
    priority: 2
    criterias:
      - value_name: load_avg_1
        more_than: 5`

	_, err := helper.YamlToConfig([]byte(yamlData))

	if err != nil {
		t.Error(formatYamlError(yamlData, err))
	}
}

func TestSameSchedMultiplePriorities(t *testing.T) {
	yamlData := `interval: 1000
schedulers:
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 5
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 2
    criterias:
      - value_name: load_avg_1
        more_than: 5`

	_, err := helper.YamlToConfig([]byte(yamlData))

	if err != nil {
		t.Error(formatYamlError(yamlData, err))
	}
}

// //////// Negative tests //////////
func TestInvalidValueName(t *testing.T) {
	yamlData := `interval: 1000
schedulers:
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg
        less_than: 1`

	_, err := helper.YamlToConfig([]byte(yamlData))

	var expErr *errs.InvalidValueNameError
	if !errors.As(err, &expErr) {
		t.Error(formatYamlError(yamlData, err))
	}
}

func TestMissingParameter(t *testing.T) { // no more_than or less_than specified
	yamlData := `interval: 1000
schedulers:
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg_1`

	_, err := helper.YamlToConfig([]byte(yamlData))

	var expErr *errs.MissingParameterError
	if !errors.As(err, &expErr) {
		t.Error(formatYamlError(yamlData, err))
	}
}

func TestConflictCriterias(t *testing.T) { // same value_name specified multiple times
	yamlData := `interval: 1000
schedulers:
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 1
      - value_name: load_avg_1
        more_than: 1`

	_, err := helper.YamlToConfig([]byte(yamlData))

	var expErr *errs.ConflictCriteriasError
	if !errors.As(err, &expErr) {
		t.Error(formatYamlError(yamlData, err))
	}
}

func TestConflictParametersBigger(t *testing.T) { // more_than > less_than
	yamlData := `interval: 1000
schedulers:
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 1
        more_than: 1.1`

	_, err := helper.YamlToConfig([]byte(yamlData))

	var expErr *errs.ConflictParametersError
	if !errors.As(err, &expErr) {
		t.Error(formatYamlError(yamlData, err))
	}
}

func TestConflictParametersEqual(t *testing.T) { // more_than == less_than
	yamlData := `interval: 1000
schedulers:
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 1
        more_than: 1`

	_, err := helper.YamlToConfig([]byte(yamlData))

	var expErr *errs.ConflictParametersError
	if !errors.As(err, &expErr) {
		t.Error(formatYamlError(yamlData, err))
	}
}

func TestPathDoesNotExist(t *testing.T) {
	yamlData := `interval: 1000
schedulers:
  - path: "invalid_path.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 1`

	_, err := helper.YamlToConfig([]byte(yamlData))

	if !errors.Is(err, os.ErrNotExist) {
		t.Error(formatYamlError(yamlData, err))
	}
}

func TestConflictPriorities(t *testing.T) { // both schedulers have same priority
	yamlData := `interval: 1000
schedulers:
  - path: "../../bytecode/fifo.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 1
  - path: "../../bytecode/lottery.bpf.c.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 1`

	_, err := helper.YamlToConfig([]byte(yamlData))

	var expErr *errs.ConflictPrioritiesError
	if !errors.As(err, &expErr) {
		t.Error(formatYamlError(yamlData, err))
	}
}
