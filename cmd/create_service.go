package cmd

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var createSvcCmd = &cobra.Command{
	Use:       "new",
	Short:     "Create a new project with the given name",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"name"},
	RunE: func(cmd *cobra.Command, args []string) error {
		input := map[string]string{
			"Name":   args[0],
			"Module": args[0],
		}

		if err := os.Mkdir(input["Name"], 0755); err != nil {
			return err
		}

		err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
			if path == "templates" {
				return nil
			}

			outPath := strings.Replace(path, "templates", input["Name"], 1)
			outPath = strings.Replace(outPath, ".tmpl", "", 1)

			if info.IsDir() {
				os.Mkdir(outPath, 0755)
				return nil
			}

			tmpl, err := template.ParseFiles(path)
			if err != nil {
				return err
			}

			file, err := os.Create(outPath)
			if err != nil {
				return err
			}

			if err := tmpl.Execute(file, input); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	},
}
