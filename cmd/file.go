package cmd

import (
	"github.com/spf13/cobra"
	"hotalert/logging"
	"hotalert/task"
	"hotalert/task/target"
	"hotalert/workload"
	"os"
	"sync"
)

// The file command executes a tasks from a single file only.
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "execute tasks from a single file",
	Run: func(cmd *cobra.Command, args []string) {
		var isRunning = true
		var fileName = args[0]
		data, err := os.ReadFile(fileName)
		if err != nil {
			logging.SugaredLogger.Fatalf("Failed to read file %s exiting!", fileName)
			return
		}
		workload, err := workload.FromYamlContent(data)
		if err != nil {
			logging.SugaredLogger.Fatalf("Failed to parse file %s exiting!", fileName)
			return
		}

		var waitGroup = sync.WaitGroup{}
		var defaultExecutor = task.NewDefaultExecutor(target.ScrapeWebTask)
		taskResultChan := defaultExecutor.Start()
		defer defaultExecutor.Shutdown()

		// Function that logs task results and marks the task executed in the waitGroup
		go func() {
			for isRunning {
				select {
				case result := <-taskResultChan:
					waitGroup.Done()
					if result.Error() != nil {
						logging.SugaredLogger.Errorf("Failed to execute task %v got: %s", result.InitialTask, result.Error())
					}
				}
			}
		}()

		// Add tasks
		waitGroup.Add(workload.GetTasksLen())
		for _, task := range workload.GetTasks() {
			defaultExecutor.AddTask(task)
		}

		// Wait for tasks to be executed
		waitGroup.Wait()

		// Turn off logging goroutine
		isRunning = false

		logging.SugaredLogger.Infof("Done")
	},
}
