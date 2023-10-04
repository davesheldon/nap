/*
Copyright © 2021 Bold City Software

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

cmd/generate.go - this is the handler for the generate command
*/
package cmd

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"github.com/kennygrant/sanitize"
// 	"github.com/spf13/cobra"
// )

// var supportedComponentTypes = []string{"request", "routine", "env", "script"}

// func containsString(slice []string, value string) bool {
// 	for _, s := range slice {
// 		if s == value {
// 			return true
// 		}
// 	}

// 	return false
// }

// // generateCmd represents the generate command
// var generateCmd = &cobra.Command{
// 	Use:   "generate <type> <target>",
// 	Short: "Generate a new request, routine, script or environment",
// 	Long:  `The generate command will add a file to the appropriate location using the template for that object type.`,
// 	Args:  cobra.MinimumNArgs(2),
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		if !containsString(supportedComponentTypes, args[0]) {
// 			return fmt.Errorf("generate requires a valid type argument. valid options: %s", strings.Join(supportedComponentTypes, ", "))
// 		}

// 		var componentType = args[0]
// 		var componentName = args[1]

// 		filePath := getFilePath(componentType, componentName)

// 		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {

// 			data, err := getTemplateData(componentType, componentName)
// 			if err != nil {
// 				return err
// 			}

// 			err = os.WriteFile(filePath, data, 0644)
// 			if err != nil {
// 				return err
// 			}

// 			fmt.Printf("created %s %s (path: %s)", componentType, componentName, filePath)
// 		} else {
// 			return fmt.Errorf("unable to create %s %s. file already exists (path: %s)", componentType, componentName, filePath)
// 		}

// 		return nil
// 	},
// }

// func getFilePath(componentType string, componentName string) string {
// 	folder := getComponentFolder(componentType)
// 	extension := getComponentExtension(componentType)

// 	fileName := sanitize.Name(componentName) + extension

// 	return filepath.Join(folder, fileName)
// }

// func getComponentFolder(componentType string) string {
// 	switch componentType {
// 	case "request":
// 		return "requests"
// 	case "routine":
// 		return "routines"
// 	case "env":
// 		return "env"
// 	case "script":
// 		return "scripts"
// 	}

// 	return ""
// }

// func getComponentExtension(componentType string) string {
// 	switch componentType {
// 	case "request":
// 		return ".yml"
// 	case "routine":
// 		return ".yml"
// 	case "env":
// 		return ".yml"
// 	case "script":
// 		return ".js"
// 	}

// 	return ""
// }

// func getTemplateData(componentType string, componentName string) ([]byte, error) {
// 	templatePath := fmt.Sprintf(".templates/%s%s", componentType, getComponentExtension(componentType))
// 	if _, err := os.Stat(templatePath); errors.Is(err, os.ErrNotExist) {
// 		return nil, fmt.Errorf("could not find project %s template (looked in %s)", componentType, templatePath)
// 	} else if err != nil {
// 		return nil, fmt.Errorf("error opening project %s template (path: %s)", componentType, templatePath)
// 	}

// 	data, err := os.ReadFile(templatePath)

// 	if err != nil {
// 		return nil, fmt.Errorf("error opening project %s template (path: %s)", componentType, templatePath)
// 	}

// 	data = []byte(strings.ReplaceAll(string(data), "${name}", componentName))

// 	return data, nil
// }

// func init() {
// 	rootCmd.AddCommand(generateCmd)
// 	rootCmd.CompletionOptions.DisableDefaultCmd = true
// }
