package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createSvcCmd = &cobra.Command{
	Use:       "new",
	Short:     "Create a new project with the given name",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"name"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating project", args)
	},
}
