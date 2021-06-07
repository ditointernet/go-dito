package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "go-dito",
	Version: "0.0.1",
}

func init() {
	migrateCmd.PersistentFlags().UintP(
		"steps", "s", 0,
		"number of steps operation. will apply all migrations if not informed or steps equal zero",
	)

	migrateCmd.AddCommand(upCmd, downCmd, seedCmd)

	rootCmd.AddCommand(createSvcCmd, migrateCmd)
}

// Execute ...
func Execute() error {
	return rootCmd.Execute()
}

type sysCommand struct {
	cmd  string
	args []string
}

func runSysCmds(dir string, cmds []*sysCommand) error {
	for _, c := range cmds {
		sysCmd := exec.Command(c.cmd, c.args...)
		sysCmd.Dir = dir
		if err := sysCmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
