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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/davesheldon/nap/naproutine"
	"github.com/robertkrimen/otto"

	"github.com/davesheldon/nap/naprequest"
	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/naprunner"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <target>",
	Short: "executes the target",
	Long:  `The run command executes the request, routine or script at the path provided.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		runConfig := newRunConfig(cmd, args)

		if runConfig.Verbose {
			runConfig.printStats()
		}

		environmentVariables, err := loadEnvironment(runConfig)
		if err != nil {
			return err
		}

		ctx, err := napcontext.New(environmentVariables)
		if err != nil {
			return err
		}

		routineResult := naprunner.RunPath(ctx, runConfig.Target)

		if routineResult.Error != nil {
			return routineResult.Error
		}

		routineResult.Print("")

		return nil
	},
}

func loadEnvironment(runConfig *RunConfig) (map[string]string, error) {
	environmentVariables := make(map[string]string)

	if len(runConfig.Environment) > 0 {
		environmentFileName := path.Join("env", runConfig.Environment+".yml")

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
	Environment string
	Verbose     bool
}

func (c *RunConfig) printStats() {
	fmt.Printf("Target: %s\n", c.Target)
	fmt.Printf("Environment: %s\n", c.Environment)
	fmt.Printf("Verbose Mode: %t\n", c.Verbose)
}

func newRunConfig(cmd *cobra.Command, args []string) *RunConfig {
	config := new(RunConfig)
	config.Target = args[1]
	config.Environment, _ = cmd.Flags().GetString("env")

	if len(config.Environment) == 0 {
		config.Environment = "default"
	}

	config.Verbose, _ = cmd.Flags().GetBool("verbose")
	return config
}

func readBodyAsString(httpResponse *http.Response) (string, error) {
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func jsonPretty(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("env", "e", "", "name of the environment variable set to use")
}
