package cmd

import (
	"github.com/spf13/cobra"
	"hotalert/logging"
)

var directoryCmd = &cobra.Command{
	Use:   "directory",
	Short: "execute each yaml file from a directory",
	Run: func(cmd *cobra.Command, args []string) {
		logging.SugaredLogger.Fatal("not implemented")
	},
}
