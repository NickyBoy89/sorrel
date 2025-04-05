package main

import (
	"github.com/spf13/cobra"
)

var configFileLocation = "sorrel-config.json"

func main() {
	mainCmd := &cobra.Command{
		Use: "sorrel",
	}

	mainCmd.Flags().StringVar(&configFileLocation, "config", "sorrel-config.json", "Location of the configuration file")

	mainCmd.AddCommand(serveCommand)

	if err := mainCmd.Execute(); err != nil {
		panic(err)
	}
}
