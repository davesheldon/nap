/*
Copyright © 2021 Dave Sheldon <dave@boldcitysoftware.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/davesheldon/nap/internal"

	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <type> <name>",
	Short: "Create a new object",
	Long:  `The create command will add a file and stub out a request skeleton for it.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || (args[0] != "request" && args[0] != "config") {
			return errors.New("create requires a valid type argument. valid options: request, config")
		}

		if len(args) < 2 {
			return errors.New("create requires a name")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var objectType = args[0]
		var objectName = args[1]

		fileName := getFileName(objectType, objectName)

		if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {

			data, err := getTemplateData(objectType, objectName)

			if err != nil {
				return err
			}

			err = os.WriteFile(fileName, data, 0644)
			if err != nil {
				return err
			}

			stubType := "request"

			if args[0] == "config" {
				stubType = "configuration file"
			}

			fmt.Printf("Created new %s stub: %s", stubType, fileName)
		} else {
			return errors.New(fmt.Sprintf("The file '%s' already exists.", fileName))
		}

		return nil
	},
}

func getFileName(objectType string, objectName string) string {
	switch objectType {
	case "request":
		if !strings.HasSuffix(objectName, ".yml") {
			objectName = objectName + ".yml"
		}
		break
	case "config":
		if !strings.HasSuffix(objectName, ".yml") {
			objectName = objectName + ".yml"
		}
		break
	}

	return sanitize.Path(objectName)
}

func getTemplateData(objectType string, objectName string) ([]byte, error) {
	switch objectType {
	case "request":
		return internal.NewNapRequest(objectName).ToYaml()
	default:
		return []byte{}, nil
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
