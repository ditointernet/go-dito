package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "go-dito",
	Version: "0.0.1",
}

func init() {
	rootCmd.AddCommand(createSvcCmd)
}

// Execute ...
func Execute() error {
	return rootCmd.Execute()
}
