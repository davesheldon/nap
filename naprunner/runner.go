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
	"fmt"
	"github.com/davesheldon/nap/naprequest"
	"github.com/davesheldon/nap/naproutine"
	"github.com/davesheldon/nap/napcontext"
	"github.com/kennygrant/sanitize"
	"gopkg.in/yaml.v2"
	"os"
	"path"
    "path/filepath"
	"strings"
	"sync"
	"time"
)

func RunPath(ctx *Context, path string) *RoutineResult {
	routine := naproutine.NewRoutine(ctx, path)
	return runRoutine(ctx, routine, nil, nil)
}

func runScript(ctx *Context, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	script := string(data)

	return runScriptInline(ctx, script)
}

func runScriptInline(ctx *Context, script string) error {
	_, err = ctx.ScriptVm.Run(script)

	return err
}

func runRequest(ctx *Context, request *naprequest.Request) *RequestResult {
	result := new(RequestResult)

	result.StartTime = time.Now()

	if len(request.PreRequestScript) > 0 {
		runScriptInline(ctx, request.PreRequestScript)
	}

	if len(request.PreRequestScriptFile) > 0 {
		runScript(ctx, request.PreRequestScriptFile)
	}

	result.HttpResponse, result.Error = executeHttp(request)

	if len(request.PostRequestScript) > 0 {
		runScriptInline(ctx, request.PostRequestScript)
	}

	if len(request.PostRequestScriptFile) > 0 {
		runScript(ctx, request.PostRequestScriptFile)
	}

	result.EndTime = time.Now()

	return result
}

func executeHttp(r *naprequest.Request) (*http.Response, error) {
	client := &http.Client{}

	var content io.Reader

	if len(r.Body) > 0 {
		content = bytes.NewBuffer([]byte(r.Body))
	} else {
		content = strings.NewReader("")
	}

	request, err := http.NewRequest(r.Verb, r.Path, content)

	if err != nil {
		return nil, err
	}

	for k, v := range r.Headers {
		request.Header.Add(k, v)
	}

	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	return resp, nil
}


func runRoutine(ctx *Context, routine *Routine, parentStep *RoutineStep, ch chan *RoutineStepResult) *RoutineResult {
	result := new(RoutineResult)
	result.StartTime = time.Now()

	waitCount := 0

	childCh := make(chan *RoutineStepResult)

	var initialCtx := routine.Context.EnvironmentVariables.Clone()

	for _, step := range routine.Steps {
		var stepResult *RoutineStepResult

		stepType, err := peekType(step.Run)

		if err != nil {
			stepResult = StepError(step, err)
		}
		else {
			if stepType == "request" {
				request, err := naprequest.LoadFromPath(step.Run, ctx)
				
				if err != nil {
					stepResult = StepError(step, err)
				} else {
					stepResult = StepRequestResult(step, runRequest(ctx, request))
				}
			}

			if step.Type == "script" {
				err = runScript(ctx, step.Run)

				if err != nil {
					stepResult = StepError(step, err)
				} else {
					routine.Context.ScriptVm.Run(string(scriptData))
					stepResult = StepScriptResult(step)
				}
				
			}

			if step.Type == "routine" {
				subroutineCtx := initialCtx.Clone()
				subroutine, err := LoadFromPath(step.Run, subroutineCtx)

				if err != nil {
					stepResult = StepError(step, err)
				} else {
					waitCount = waitCount + 1

					go runRoutine(subroutineCtx, subroutine, step, childCh)
					// we'll get the results after the loop finishes
					continue
				}
			}
		}

		result.StepResults = append(result.StepResults, stepResult)
	}

	for i := 0; i < waitCount; i++ {
		deferredResult := <-childCh
		result.StepResults = append(result.StepResults, deferredResult)
	}

	if ch != nil {
		ch <- StepSubroutineResult(parentStep, result)
	}

	result.EndTime = time.Now()

	return result
}

func peekType(path string) (string, error) {
	if filepath.Ext(path) == ".js" {
		return "script"
	}

	var yamlMap := map[string]string{}

	data, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	err := yaml.Unmarshal(data, &yamlMap)

	if val, ok := dict["foo"]; ok {
		return val, nil
	}
	else {
		return "", errors.New("type not found")
	}
}
