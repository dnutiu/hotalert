package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "hotalert",
	Args:  cobra.ExactArgs(1),
	Short: "Hotalert is a command line tool that for task execution and configuration.",
	Long: `Hotalert is a command line tool that for task execution and configuration. Tasks and alerts are defined 
in yaml files and the program parses the files, executes the tasks and emits alerts when the tasks conditions are met. `,
}

func init() {
	RootCmd.AddCommand(fileCmd)
	RootCmd.AddCommand(directoryCmd)
}
