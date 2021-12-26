/*
Copyright © 2021 Dave Sheldon <dave@boldcitysoftware.com>
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
	"path/filepath"
	"strings"

	"github.com/davesheldon/nap/internal"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <target>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("run requires a valid request name.")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		runConfig := newRunConfig(cmd, args)

		if runConfig.Verbose {
			runConfig.printStats()
		}

		configMap, err := loadConfigMap(runConfig)
		if err != nil {
			return err
		}

		file, err := os.Open(runConfig.TargetFileName)
		if err != nil {
			return err
		}

		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			if runConfig.Verbose {
				fmt.Printf("Walking directory: %s\n", runConfig.TargetFileName)
			}

			err = filepath.Walk(runConfig.TargetFileName, func(path string, f os.FileInfo, _ error) error {
				if !f.IsDir() && filepath.Ext(path) == ".yml" {
					err = executeFile(path, configMap, cmd)
					if err != nil {
						return err
					}
				}

				return nil
			})

			return err
		} else {
			return executeFile(runConfig.TargetFileName, configMap, cmd)
		}
	},
}

func loadConfigMap(runConfig *RunConfig) (map[string]string, error) {
	configMap := make(map[string]string)

	if len(runConfig.ConfigFileName) > 0 {
		if !strings.HasSuffix(runConfig.ConfigFileName, ".yml") {
			runConfig.ConfigFileName = runConfig.ConfigFileName + ".yml"
		}

		if _, err := os.Stat(runConfig.ConfigFileName); errors.Is(err, os.ErrNotExist) {
			return configMap, errors.New(fmt.Sprintf("config file name '%s' not found.", runConfig.ConfigFileName))
		} else if err != nil {
			return configMap, err
		}

		configData, err := os.ReadFile(runConfig.ConfigFileName)
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
	TargetFileName string
	ConfigFileName string
	Verbose        bool
}

func (c *RunConfig) printStats() {
	fmt.Printf("Target File Name: %s\n", c.TargetFileName)
	fmt.Printf("Config File Name: %s\n", c.ConfigFileName)
	fmt.Printf("Verbose Mode: %t\n", c.Verbose)
}

func newRunConfig(cmd *cobra.Command, args []string) *RunConfig {
	config := new(RunConfig)
	config.TargetFileName = args[0]
	config.ConfigFileName, _ = cmd.Flags().GetString("config")
	config.Verbose, _ = cmd.Flags().GetBool("verbose")
	return config
}

func executeFile(fileName string, config map[string]string, cmd *cobra.Command) error {
	runnableData, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	runnableTemplate := string(runnableData)

	for k, v := range config {
		variable := fmt.Sprintf("${%s}", k)
		runnableTemplate = strings.ReplaceAll(runnableTemplate, variable, v)
	}

	runnableData = []byte(runnableTemplate)

	runnable, err := internal.ParseNapRunnable(runnableData)

	if err != nil {
		return err
	}

	fmt.Printf("- %s: ", fileName)

	return executeRunnable(runnable, cmd)
}

func executeRunnable(runnable *internal.NapRunnable, cmd *cobra.Command) error {
	verbose, _ := cmd.Flags().GetBool("verbose")

	if verbose {
		runnable.PrintStats()
	}

	result := runnable.GetResult()

	if result.Error != nil {
		return result.Error
	}

	if verbose {
		fmt.Printf("Response Status: %s (Content Length: %d bytes)\n", result.HttpResponse.Status, result.HttpResponse.ContentLength)
	} else {
		fmt.Println(result.HttpResponse.Status)
	}

	if verbose {
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

	runCmd.Flags().StringP("config", "c", "", "Path to an external configuration file.")
}
