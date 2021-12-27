/*
Package cmd
Copyright © 2021 Bold City Software <dave@boldcitysoftware.com>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kennygrant/sanitize"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/davesheldon/nap/internal"

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
			targetFileName := path.Join("requests", sanitize.BaseName(runConfig.Target))
			targetFile, err := os.Open(targetFileName)

			if err != nil {
				return err
			}

			defer targetFile.Close()

			result := runRequest(runConfig, targetFileName, environmentVariables)

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
			return errors.New("running routines is not yet implemented")
		}

		return nil
	},
}

func loadEnvironment(runConfig *RunConfig) (map[string]string, error) {
	configMap := make(map[string]string)

	if len(runConfig.Environment) > 0 {
		if !strings.HasSuffix(runConfig.Environment, ".yml") {
			runConfig.Environment = runConfig.Environment + ".yml"
		}

		if _, err := os.Stat(runConfig.Environment); errors.Is(err, os.ErrNotExist) {
			return configMap, errors.New(fmt.Sprintf("config file name '%s' not found.", runConfig.Environment))
		} else if err != nil {
			return configMap, err
		}

		configData, err := os.ReadFile(runConfig.Environment)
		if err != nil {
			return configMap, err
		}

		err = yaml.Unmarshal(configData, &configMap)
		if err != nil {
			return configMap, err
		}
	}

	return configMap, nil
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

func runRequest(runConfig *RunConfig, fileName string, environmentVariables map[string]string) *internal.NapRequestResult {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return internal.NapRequestResultError(err)
	}

	dataAsString := string(data)

	for k, v := range environmentVariables {
		variable := fmt.Sprintf("${%s}", k)
		dataAsString = strings.ReplaceAll(dataAsString, variable, v)
	}

	data = []byte(dataAsString)

	request, err := internal.ParseNapRequest(data)

	if err != nil {
		return internal.NapRequestResultError(err)
	}

	return executeRunnable(request, runConfig)
}

func executeRunnable(request *internal.NapRequest, runConfig *RunConfig) *internal.NapRequestResult {
	if runConfig.Verbose {
		request.PrintStats()
	}

	return request.GetResult()
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
