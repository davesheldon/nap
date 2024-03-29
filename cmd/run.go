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
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/napenv"
	"github.com/davesheldon/nap/naprunner"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <target>",
	Short: "Execute a request, routine or script",
	Long:  `The run command executes a request, routine or script at the path provided.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		runConfig := newRunConfig(cmd, args)

		environmentVariables, err := loadEnvironment(runConfig)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		var wg sync.WaitGroup
		napCtx := napcontext.New(".", runConfig.Environments, environmentVariables, &wg, runConfig.Quiet)

		routineResult := naprunner.RunPath(napCtx, runConfig.Target)

		napCtx.Complete()

		end := time.Now()

		if runConfig.Verbose {
			if len(routineResult.StepResults) == 1 && routineResult.StepResults[0].SubroutineResult != nil {
				routineResult.StepResults[0].SubroutineResult.Print("", napCtx)
			} else {
				routineResult.Print("", napCtx)
			}
		} else {
			for _, err := range routineResult.Errors {
				fmt.Printf("[ERROR] %s\n", err.Error())
			}
		}

		runStats := routineResult.GetRunStats()

		runTypes := make([]string, 0, len(runStats.StatsByType))

		for k := range runStats.StatsByType {
			runTypes = append(runTypes, k)
		}
		sort.Strings(runTypes)

		var statsPerTypeOutput string

		for _, runType := range runTypes {
			typeStats := runStats.StatsByType[runType]
			statsPerTypeOutput += fmt.Sprintf("%s:\t%d/%d\n", runType, typeStats.Passing, typeStats.Total)
		}

		statsPerTypeOutput += fmt.Sprintf("Total:\t\t%d/%d", runStats.Totals.Passing, runStats.Totals.Total)

		if runStats.Totals.Total == runStats.Totals.Passing {
			fmt.Printf("\n%s\n\nSUCCESS! Run finished in %dms.\n", statsPerTypeOutput, end.Sub(start).Milliseconds())
			return nil
		} else {
			cmd.SilenceUsage = true
			fmt.Printf("\n%s\n\n", statsPerTypeOutput)
			return fmt.Errorf("Run failed after %dms.", end.Sub(start).Milliseconds())
		}

	},
}

func loadEnvironment(runConfig *RunConfig) (map[string]string, error) {
	environmentVariables := make(map[string]string)

	for _, environmentFileNameOriginal := range runConfig.Environments {
		result, err := napenv.AddEnvironmentFromPath(runConfig.TargetDir, environmentFileNameOriginal, environmentVariables)
		if err != nil {
			return environmentVariables, err
		}

		environmentVariables = result
	}

	for k, v := range runConfig.Variables {
		environmentVariables[k] = v
	}

	return environmentVariables, nil
}

type RunConfig struct {
	Target       string
	TargetDir    string
	TargetName   string
	Environments []string
	Variables    map[string]string
	Verbose      bool
	Quiet        bool
}

func newRunConfig(cmd *cobra.Command, args []string) *RunConfig {
	config := new(RunConfig)
	config.Target = args[0]
	config.TargetDir = filepath.Dir(config.Target)
	config.TargetName = path.Base(config.Target)
	config.Environments, _ = cmd.Flags().GetStringArray("env")
	config.Verbose, _ = cmd.Flags().GetBool("verbose")
	config.Variables = make(map[string]string)
	config.Quiet, _ = cmd.Flags().GetBool("quiet")

	params, _ := cmd.Flags().GetStringArray("param")

	for _, p := range params {
		keyVal := strings.Split(p, "=")
		if len(keyVal) == 2 {
			config.Variables[keyVal[0]] = keyVal[1]
		}
	}

	return config
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringArrayP("env", "e", []string{}, "add environment variables from a file `path`")
	runCmd.Flags().StringArrayP("param", "p", []string{}, "add a single variable to the run as a `<name>=<value>` pair")
	runCmd.Flags().BoolP("quiet", "q", false, "suppress output until the end")
}
