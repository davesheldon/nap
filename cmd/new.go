/*
Package cmd
Copyright © 2021 Bold City Software <dave@boldcitysoftware.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/kennygrant/sanitize"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new project",
	Long:  `The new command creates a new Nap project in a folder named after project name`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		projectPath := sanitize.BaseName(projectName)

		if err := ensureDirectoryExists(projectPath); err != nil {
			return err
		}

		subdirectoriesToCreate := []string{"requests", "routines", "env", ".template"}

		for _, subdirectory := range subdirectoriesToCreate {
			fullPath := path.Join(projectPath, subdirectory)

			if err := ensureDirectoryExists(fullPath); err != nil {
				return err
			}
		}

		requestTemplateData := []byte(`name: ${name}
path: https://cat-fact.herokuapp.com/facts
verb: GET
body:
headers:
    - Accept: application/json
	- Content-Type: application/json
`)
		if err := tryWriteFileData(path.Join(projectPath, ".template", "request.yml"), requestTemplateData); err != nil {
			return err
		}

		routineTemplateData := []byte(`name: ${name}
run:
  - type: request
    name: 
`)
		if err := tryWriteFileData(path.Join(projectPath, ".template", "routine.yml"), routineTemplateData); err != nil {
			return err
		}

		defaultEnvData := []byte(``)
		if err := tryWriteFileData(path.Join(projectPath, "env", "default.yml"), defaultEnvData); err != nil {
			return err
		}

		firstRequestData := []byte(`name: request-1
path: https://cat-fact.herokuapp.com/facts
verb: GET
body:
headers:
    - Accept: application/json
	- Content-Type: application/json
`)
		if err := tryWriteFileData(path.Join(projectPath, "requests", "request-1.yml"), firstRequestData); err != nil {
			return err
		}

		firstRoutineData := []byte(`name: routine-1
run:
    - type: request 
      name: my-request
      expectStatusCode: 200`)
		if err := tryWriteFileData(path.Join(projectPath, "routines", "routine-1.yml"), firstRoutineData); err != nil {
			return err
		}

		fmt.Printf("Project '%s' created.\n", projectPath)

		return nil
	},
}

func tryWriteFileData(filePath string, data []byte) error {
	fmt.Println("in tryWriteFileData")

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(filePath, data, os.ModePerm)
	}

	return errors.New(fmt.Sprintf("could not create %s. either a file already exists or the path isn't accessible", filePath))
}

func ensureDirectoryExists(directoryPath string) error {
	projectPathInfo, err := os.Stat(directoryPath)

	if errors.Is(err, os.ErrNotExist) {
		os.Mkdir(directoryPath, os.ModePerm)
	} else if !projectPathInfo.IsDir() {
		return errors.New(fmt.Sprintf("path exists and is not a directory (path: %s)", directoryPath))
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
