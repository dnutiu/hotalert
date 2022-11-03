package main

import (
	"hotalert/cmd"
	"hotalert/logging"
)

func main() {
	logging.InitLoggingWithParams("info", "console")
	err := cmd.RootCmd.Execute()
	if err != nil {
		logging.SugaredLogger.Fatal(err)
	}
}
