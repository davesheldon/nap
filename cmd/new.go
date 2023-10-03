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

cmd/new.go - this is the handler for the new command
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/kennygrant/sanitize"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new project",
	Long:  `The new command creates a new Nap project in a folder named after project name`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("new requires a project name")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		projectPath := sanitize.BaseName(projectName)

		if err := ensureDirectoryExists(projectPath); err != nil {
			return err
		}

		subdirectoriesToCreate := []string{".templates", "env", "requests", "routines", "scripts"}

		for _, subdirectory := range subdirectoriesToCreate {
			fullPath := filepath.Join(projectPath, subdirectory)

			if err := ensureDirectoryExists(fullPath); err != nil {
				return err
			}
		}

		requestTemplateData := []byte(`kind: request
name: Request Name
path: https://cat-fact.herokuapp.com/facts
verb: GET
body:
headers:
  Accept: application/json
  Content-Type: application/json
`)
		if err := tryWriteFileData(path.Join(projectPath, ".templates", "request.yml"), requestTemplateData); err != nil {
			return err
		}

		routineTemplateData := []byte(`kind: routine
name: Routine Name
steps:
  - run: routine.yml
`)
		if err := tryWriteFileData(path.Join(projectPath, ".templates", "routine.yml"), routineTemplateData); err != nil {
			return err
		}

		defaultEnvData := []byte(``)
		if err := tryWriteFileData(path.Join(projectPath, "env", "default.yml"), defaultEnvData); err != nil {
			return err
		}

		firstRequestData := []byte(`kind: request
name: Request 1
path: https://cat-fact.herokuapp.com/facts
verb: GET
body:
headers:
  Accept: application/json
  Content-Type: application/json
`)
		if err := tryWriteFileData(path.Join(projectPath, "requests", "request-1.yml"), firstRequestData); err != nil {
			return err
		}

		firstRoutineData := []byte(`kind: routine
name: Routine 1
steps:
  - run: ../requests/request-1.yml
`)
		if err := tryWriteFileData(path.Join(projectPath, "routines", "routine-1.yml"), firstRoutineData); err != nil {
			return err
		}

		fmt.Printf("Project '%s' created.\n", projectPath)

		firstScriptData := []byte(`console.log("Hello, World!")
		`)

		if err := tryWriteFileData(path.Join(projectPath, "scripts", "script-1.js"), firstScriptData); err != nil {
			return err
		}

		return nil
	},
}

func tryWriteFileData(filePath string, data []byte) error {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(filePath, data, os.ModePerm)
	}

	return fmt.Errorf("could not create %s. either a file already exists or the path isn't accessible", filePath)
}

func ensureDirectoryExists(directoryPath string) error {
	projectPathInfo, err := os.Stat(directoryPath)

	if errors.Is(err, os.ErrNotExist) {
		os.Mkdir(directoryPath, os.ModePerm)
	} else if !projectPathInfo.IsDir() {
		return fmt.Errorf("path exists and is not a directory (path: %s)", directoryPath)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
