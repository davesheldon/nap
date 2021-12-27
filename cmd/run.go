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
	"github.com/davesheldon/nap/naproutine"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/davesheldon/nap/naprequest"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <type> <target>",
	Short: "execute a request or routine",
	Long:  `The run command executes a request or routine using the name provided.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || (args[0] != "request" && args[0] != "env" && args[0] != "routine") {
			return errors.New("run requires a valid type argument. valid options: request, routine")
		}

		if len(args) < 2 {
			return errors.New(fmt.Sprintf("run requires a %s name", args[0]))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		runConfig := newRunConfig(cmd, args)

		if runConfig.Verbose {
			runConfig.printStats()
		}

		environmentVariables, err := loadEnvironment(runConfig)
		if err != nil {
			return err
		}

		if runConfig.TargetType == "request" {
			result := runRequest(runConfig, environmentVariables)

			if runConfig.Verbose {
				fmt.Printf("Response Status: %s (Content Length: %d bytes)\n", result.HttpResponse.Status, result.HttpResponse.ContentLength)
			} else {
				fmt.Println(result.HttpResponse.Status)
			}

			if runConfig.Verbose {
				output, err := readBodyAsString(result.HttpResponse)
				if err != nil {
					return err
				}

				if strings.Contains(result.HttpResponse.Header.Get("Content-Type"), "json") {
					output, err = jsonPretty(output)

					if err != nil {
						return err
					}
				}

				print(output)
			}

			defer result.HttpResponse.Body.Close()

			return nil
		}

		if runConfig.TargetType == "routine" {
			routine, err := naproutine.LoadByName(runConfig.Target, environmentVariables)
			if err != nil {
				return err
			}

			routineResult := routine.Run(environmentVariables, nil, nil)

			if routineResult.Error != nil {
				return routineResult.Error
			}

			routineResult.Print("")
		}

		return nil
	},
}

func loadEnvironment(runConfig *RunConfig) (map[string]string, error) {
	environmentVariables := make(map[string]string)

	if len(runConfig.Environment) > 0 {
		environmentFileName := path.Join("env", runConfig.Environment+".yml")

		if _, err := os.Stat(environmentFileName); errors.Is(err, os.ErrNotExist) {
			return environmentVariables, errors.New(fmt.Sprintf("environment '%s' not found.", runConfig.Environment))
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
	TargetType  string
	Target      string
	Environment string
	Verbose     bool
}

func (c *RunConfig) printStats() {
	fmt.Printf("Target Type: %s\n", c.TargetType)
	fmt.Printf("Target: %s\n", c.Target)
	fmt.Printf("Environment: %s\n", c.Environment)
	fmt.Printf("Verbose Mode: %t\n", c.Verbose)
}

func newRunConfig(cmd *cobra.Command, args []string) *RunConfig {
	config := new(RunConfig)
	config.TargetType = args[0]
	config.Target = args[1]
	config.Environment, _ = cmd.Flags().GetString("env")

	if len(config.Environment) == 0 {
		config.Environment = "default"
	}

	config.Verbose, _ = cmd.Flags().GetBool("verbose")
	return config
}

func runRequest(runConfig *RunConfig, environmentVariables map[string]string) *naprequest.RequestResult {
	request, err := naprequest.LoadByName(runConfig.Target, environmentVariables)

	if err != nil {
		return naprequest.ResultError(err)
	}

	if runConfig.Verbose {
		request.PrintStats()
	}

	return request.Run()
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
