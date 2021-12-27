/*
Package cmd
Copyright © 2021 Bold City Software <dave@boldcitysoftware.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate <type> <name>",
	Short: "Generate a new object",
	Long:  `The generate command will add a file to the appropriate location using the template for that object type.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || (args[0] != "request" && args[0] != "env" && args[0] != "routine") {
			return errors.New("generate requires a valid type argument. Valid options: request, routine, env")
		}

		if len(args) < 2 {
			return errors.New("generate requires a name")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var objectType = args[0]
		var objectName = args[1]

		filePath := getFilePath(objectType, objectName)

		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {

			data, err := getTemplateData(objectType, objectName)
			if err != nil {
				return err
			}

			err = os.WriteFile(filePath, data, 0644)
			if err != nil {
				return err
			}

			fmt.Printf("created %s %s (path: %s)", objectType, objectName, filePath)
		} else {
			return errors.New(fmt.Sprintf("unable to create %s %s. file already exists (path: %s)", objectType, objectName, filePath))
		}

		return nil
	},
}

func getFilePath(objectType string, objectName string) string {
	var folder string

	fileName := sanitize.Name(objectName) + ".yml"

	switch objectType {
	case "request":
		folder = "requests"
		break
	case "routine":
		folder = "routines"
		break
	case "env":
		folder = "environment"
		break
	}

	return path.Join(folder, fileName)
}

func getTemplateData(objectType string, objectName string) ([]byte, error) {
	templatePath := fmt.Sprintf(".templates/%s", objectType)
	if _, err := os.Stat(templatePath); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New(fmt.Sprintf("Could not find project %s template (looked in %s)", objectType, templatePath))
	} else if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening project %s template (path: %s)", objectType, templatePath))
	}

	data, err := os.ReadFile(templatePath)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening project %s template (path: %s)", objectType, templatePath))
	}

	data = []byte(strings.ReplaceAll(string(data), "${name}", objectName))

	return data, nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
