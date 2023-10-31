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

runner.go - this file contains logic for running requests, routines and scripts
*/
package naprunner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/naprequest"
	"github.com/davesheldon/nap/naproutine"
	"github.com/davesheldon/nap/naputil"
	"gopkg.in/yaml.v2"
)

func runRoutine(ctx *napcontext.Context, routine *naproutine.Routine, parentStep *naproutine.RoutineStep, ch chan *naproutine.RoutineStepResult) *naproutine.RoutineResult {
	result := new(naproutine.RoutineResult)
	result.Routine = routine
	result.StartTime = time.Now()

	waitCount := 0

	childCh := make(chan *naproutine.RoutineStepResult)

	var progress *napcontext.Progress

	if ch != nil {
		progress = ctx.ProgressStart(routine.Name, int64(len(routine.Steps)))
	}

	for _, step := range routine.Steps {
		step.SetupContext(ctx)
		var stepResult *naproutine.RoutineStepResult
		stepResult = nil

		stepPath := filepath.Join(ctx.WorkingDirectory, step.Run)

		if exists, _ := naputil.FileExists(stepPath); !exists {
			stepResult = naproutine.StepError(step, fmt.Errorf("file doesn't exist: %s", stepPath))
			result.StepResults = append(result.StepResults, stepResult)
			if ch != nil {
				ctx.ProgressCancel(progress)
			}
			break
		}

		stepType, err := peekType(stepPath, ctx)

		if err != nil {
			stepResult = naproutine.StepError(step, err)
			result.StepResults = append(result.StepResults, stepResult)
			if ch != nil {
				ctx.ProgressCancel(progress)
			}
			break
		}

		iterations, err := step.GetIterations(ctx)

		if err != nil {
			stepResult = naproutine.StepError(step, err)
			result.StepResults = append(result.StepResults, stepResult)
			if ch != nil {
				ctx.ProgressCancel(progress)
			}
			break
		}

		for _, iterationCtx := range iterations {

			if stepType == "request" {
				request, err := naprequest.LoadFromPath(stepPath, iterationCtx)

				if err != nil {
					stepResult = naproutine.StepError(step, err)
					result.StepResults = append(result.StepResults, stepResult)
					break
				} else {
					stepResult = naproutine.StepRequestResult(step, runRequest(iterationCtx, stepPath, request))
				}
			}

			if stepType == "script" {
				stepResult = naproutine.StepScriptResult(step, runScript(iterationCtx, stepPath))
			}

			if stepType == "routine" {
				subroutineCtx := iterationCtx.Clone(filepath.Dir(stepPath))
				subroutine, err := naproutine.LoadFromPath(stepPath, subroutineCtx)

				if err != nil {
					stepResult = naproutine.StepError(step, err)
				} else {
					waitCount = waitCount + 1

					go runRoutine(subroutineCtx, subroutine, step, childCh)
					// we'll get the results after the loop finishes
					continue
				}
			}

			if stepResult == nil {
				stepResult = naproutine.StepError(step, fmt.Errorf("could not run path: %s", stepPath))
				if ch != nil {
					ctx.ProgressCancel(progress)
				}
				break
			}

			result.StepResults = append(result.StepResults, stepResult)
			if ch != nil {
				ctx.ProgressIncrement(progress)
			}
		}
	}

	for i := 0; i < waitCount; i++ {
		deferredResult := <-childCh
		result.StepResults = append(result.StepResults, deferredResult)
		if ch != nil {
			ctx.ProgressIncrement(progress)
		}
	}

	if ch != nil {
		ch <- naproutine.StepSubroutineResult(parentStep, result)
	}

	result.EndTime = time.Now()

	return result
}

func populateErrors(result *naproutine.RoutineResult) {
	for _, v := range result.StepResults {
		if v.SubroutineResult != nil {
			populateErrors(v.SubroutineResult)
			result.Errors = append(result.Errors, v.SubroutineResult.Errors...)
			continue
		}

		if len(v.Errors) > 0 {
			result.Errors = append(result.Errors, v.Errors...)
			continue
		}

		if v.RequestResult != nil && v.RequestResult.Error != nil {
			result.Errors = append(result.Errors, errors.New(fmt.Sprintf("%s: %s", v.RequestResult.Request.Name, v.RequestResult.Error)))
			continue
		}

		if v.ScriptResult != nil && v.ScriptResult.Error != nil {
			result.Errors = append(result.Errors, errors.New(fmt.Sprintf("%s: %s", v.Step.Run, v.ScriptResult.Error)))
			continue
		}
	}
}

func peekType(path string, ctx *napcontext.Context) (string, error) {
	if filepath.Ext(path) == ".js" {
		return "script", nil
	}

	var yamlMap map[string]interface{}

	data, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("type of file unclear: %s (cannot read file: %s)", path, err.Error())
	}

	dataAsString := string(data)

	for k, v := range ctx.EnvironmentVariables {
		variable := fmt.Sprintf("${%s}", k)
		dataAsString = strings.ReplaceAll(dataAsString, variable, v)
	}

	data = []byte(dataAsString)

	err = yaml.Unmarshal(data, &yamlMap)

	if err != nil {
		return "", fmt.Errorf("invalid YAML file: %s (cannot unmarshal: %s)", path, err.Error())
	}

	if val, ok := yamlMap["kind"]; ok {
		sKind, ok := val.(string)
		if ok {
			return sKind, nil
		} else {
			return "", fmt.Errorf("type of file unclear: %s", path)
		}
	} else {
		return "", fmt.Errorf("type of file unclear: %s", path)
	}
}
