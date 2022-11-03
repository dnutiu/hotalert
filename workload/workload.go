package workload

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"hotalert/alert"
	"hotalert/logging"
	"hotalert/task"
	"time"
)

// Workload represents a workload for the HotAlert program.
type Workload struct {
	tasksList  []*task.Task
	alerterMap map[string]alert.Alerter
}

// NewWorkload returns a new Workload given the workload data.
func NewWorkload(workloadData map[string]any) (*Workload, error) {
	var workload Workload
	var err error

	// Two important keys from here on: alert and tasks
	err = workload.buildAlerterMap(workloadData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed build alert contents: %s", err))
	}
	err = workload.buildTasksArray(workloadData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to build tasks contents: %s", err))
	}

	return &workload, nil
}

// GetTasks returns the tasks assigned to the workload.
func (p *Workload) GetTasks() []*task.Task {
	return p.tasksList
}

// GetTasksLen returns the number of the tasks assigned to the workload.
func (p *Workload) GetTasksLen() int {
	return len(p.tasksList)
}

// buildAlerterMap parses the alert section from the given workload data and creates alerter components.
// On failure, it returns an error.
func (p *Workload) buildAlerterMap(workloadData map[string]any) error {
	p.alerterMap = make(map[string]alert.Alerter)

	alertContents, ok := workloadData["alerts"]
	if !ok {
		return errors.New("key 'alerts' does not exists in workload")
	}
	alertContentsMap, ok := alertContents.(map[string]any)
	if !ok {
		return errors.New("key 'alert' is not a map type")
	}

	// Parse the map and build alerts on the way.
	for key, values := range alertContentsMap {
		valuesMap, ok := values.(map[string]any)
		if !ok {
			return errors.New(fmt.Sprintf("alert section '%s' does no contain a map", key))
		}

		alerter, err := alert.NewAlerter(key, valuesMap)
		if err != nil {
			return err
		}
		if p.alerterMap[key] != nil {
			return errors.New(fmt.Sprintf("alert section '%s' is a duplicate", key))
		}
		p.alerterMap[key] = alerter
	}

	return nil
}

// buildTasksArray parses the tasks section from the given workload data and creates task components.
// On failure, it returns an error.
func (p *Workload) buildTasksArray(workloadData map[string]any) error {
	p.tasksList = make([]*task.Task, 0, 10)

	// Figure out tasks types safely.
	taskContents, ok := workloadData["tasks"]
	if !ok {
		return errors.New("key 'tasks' does not exists in workload")
	}

	taskContentsArray, ok := taskContents.([]any)
	if !ok {
		return errors.New("key 'tasks' is not an array")
	}

	// Iterate through task array
	for i, taskEntryRaw := range taskContentsArray {
		taskEntry, ok := taskEntryRaw.(map[string]any)
		if !ok {
			logging.SugaredLogger.Errorf("error parsing entry %d in tasks array: not a valid map type", i)
			continue
		}

		taskOptions, ok := taskEntry["options"].(map[string]any)
		if !ok {
			logging.SugaredLogger.Errorf("error parsing entry %d in tasks array: options is not a valid map type", i)
			continue
		}

		// Build task
		tempTask := task.NewTask(taskOptions, alert.DummyAlerter{})

		// Timeout (optional)
		taskTimeout, ok := taskEntry["timeout"].(int)
		if ok {
			tempTask.Timeout = time.Duration(taskTimeout) * time.Second
		}

		// Alerter
		taskAlerter, ok := taskEntry["alerter"].(string)
		if ok {
			alerter, ok := p.alerterMap[taskAlerter]
			if !ok {
				logging.SugaredLogger.Errorf("error parsing entry %d in tasks array: invalid alerter", i)
				continue
			}
			tempTask.Alerter = alerter
		} else {
			logging.SugaredLogger.Errorf("error parsing entry %d in tasks array: invalid alerter", i)
			continue
		}

		p.tasksList = append(p.tasksList, tempTask)
	}

	if len(p.tasksList) == 0 {
		return errors.New("tasks list is empty or parsing has failed")
	}

	return nil
}

// FromYamlContent returns a new Workload given a yaml workload data definition.
func FromYamlContent(contents []byte) (*Workload, error) {
	var workloadData map[string]any

	err := yaml.Unmarshal(contents, &workloadData)
	if err != nil {
		newError := errors.New(fmt.Sprintf("failed to unmarshal yaml contents %s", err))
		return nil, newError
	}

	return NewWorkload(workloadData)
}
