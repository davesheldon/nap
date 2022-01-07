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

cmd/run.go - this is the handler for the run command
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/naprunner"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <target>",
	Short: "Execute a request, routine or script",
	Long:  `The run command executes a request, routine or script at the path provided.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		runConfig := newRunConfig(cmd, args)

		environmentVariables, err := loadEnvironment(runConfig)
		if err != nil {
			return err
		}

		ctx := napcontext.New(environmentVariables)

		routineResult := naprunner.RunPath(ctx, runConfig.Target)

		if routineResult.Error != nil {
			return routineResult.Error
		} else {
			fmt.Println("success")
		}

		return nil
	},
}

func isDirectory(dir string) bool {
	info, err := os.Stat(dir)

	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return info.IsDir()
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func loadEnvironment(runConfig *RunConfig) (map[string]string, error) {
	environmentVariables := make(map[string]string)

	environmentFileName := runConfig.Environment

	if path.Ext(environmentFileName) != ".yml" {
		environmentFileName = environmentFileName + ".yml"
	}

	if !fileExists(environmentFileName) {
		// try and find it relative to the target path
		environmentFileName = path.Join(runConfig.TargetDir, "..", "env", environmentFileName)
	}

	if len(runConfig.Environment) > 0 {
		if _, err := os.Stat(environmentFileName); errors.Is(err, os.ErrNotExist) {
			return environmentVariables, fmt.Errorf("environment '%s' not found", runConfig.Environment)
		} else if err != nil {
			return environmentVariables, err
		}

		configData, err := os.ReadFile(environmentFileName)
		if err != nil {
			return environmentVariables, err
		}

		err = yaml.Unmarshal(configData, &environmentVariables)
		if err != nil {
			return environmentVariables, err
		}
	}

	return environmentVariables, nil
}

type RunConfig struct {
	Target      string
	TargetDir   string
	Environment string
	Verbose     bool
}

func newRunConfig(cmd *cobra.Command, args []string) *RunConfig {
	config := new(RunConfig)
	config.Target = args[0]
	config.TargetDir = path.Dir(config.Target)
	config.Environment, _ = cmd.Flags().GetString("env")
	config.Verbose, _ = cmd.Flags().GetBool("verbose")

	return config
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("env", "e", "", "name of the environment variable set to use")
}
