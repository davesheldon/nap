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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/naprequest"
	"github.com/davesheldon/nap/naproutine"
	"github.com/robertkrimen/otto"
	"gopkg.in/yaml.v2"
)

func RunPath(ctx *napcontext.Context, path string) *naproutine.RoutineResult {
	fmt.Printf("Run path: %s\n", path)
	routine := naproutine.NewRoutine(ctx, path)
	return runRoutine(ctx, routine, nil, nil)
}

func runScript(ctx *napcontext.Context, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	script := string(data)

	fmt.Printf("Running script: %s", path)

	return runScriptInline(ctx, script)
}

func runScriptInline(ctx *napcontext.Context, script string) error {
	result, err := ctx.ScriptVm.Run(script)

	if err == nil {
		fmt.Println("result")
		fmt.Println(result.ToString())
	} else {
		fmt.Println("error")
		fmt.Println(err.Error())
	}

	return err
}

func runRequest(ctx *napcontext.Context, request *naprequest.Request) *naprequest.RequestResult {
	fmt.Println("request")
	result := new(naprequest.RequestResult)

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

func runRoutine(ctx *napcontext.Context, routine *naproutine.Routine, parentStep *naproutine.RoutineStep, ch chan *naproutine.RoutineStepResult) *naproutine.RoutineResult {
	fmt.Println("routine")
	result := new(naproutine.RoutineResult)
	result.StartTime = time.Now()

	waitCount := 0

	childCh := make(chan *naproutine.RoutineStepResult)

	err := setupVm(ctx)
	if err != nil {
		result.EndTime = time.Now()
		result.Error = err
		return result
	}

	for _, step := range routine.Steps {
		var stepResult *naproutine.RoutineStepResult
		stepResult = nil

		stepType, err := peekType(step.Run)

		fmt.Printf("stepType: %s\n", stepType)

		if err != nil {
			stepResult = naproutine.StepError(step, err)
		} else {
			if stepType == "request" {
				request, err := naprequest.LoadFromPath(step.Run, ctx)

				if err != nil {
					stepResult = naproutine.StepError(step, err)
				} else {
					stepResult = naproutine.StepRequestResult(step, runRequest(ctx, request))
				}
			}

			if stepType == "script" {
				err = runScript(ctx, step.Run)

				if err != nil {
					stepResult = naproutine.StepError(step, err)
				} else {
					stepResult = naproutine.StepScriptResult(step)
				}
			}

			if stepType == "routine" {
				subroutineCtx := ctx.Clone()
				subroutine, err := naproutine.LoadFromPath(step.Run, subroutineCtx)

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
				stepResult = naproutine.StepError(step, fmt.Errorf("could not run path: %s", step.Run))
			}
		}

		result.StepResults = append(result.StepResults, stepResult)
	}

	for i := 0; i < waitCount; i++ {
		deferredResult := <-childCh
		result.StepResults = append(result.StepResults, deferredResult)
	}

	if ch != nil {
		ch <- naproutine.StepSubroutineResult(parentStep, result)
	}

	result.EndTime = time.Now()

	for _, v := range result.StepResults {
		if v.Error != nil {
			result.Error = v.Error
			break
		}
	}

	return result
}

func peekType(path string) (string, error) {
	if filepath.Ext(path) == ".js" {
		return "script", nil
	}

	yamlMap := map[string]string{}

	data, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("type of file unclear: %s", path)
	}

	err = yaml.Unmarshal(data, &yamlMap)

	if err != nil {
		return "", fmt.Errorf("type of file unclear: %s", path)
	}

	if val, ok := yamlMap["type"]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("type of file unclear: %s", path)
	}
}

func setupVm(ctx *napcontext.Context) error {

	err := ctx.ScriptVm.Set("napRun", func(call otto.FunctionCall) otto.Value {

		path := call.Argument(0).String()
		RunPath(ctx, path)

		// todo: process result

		return otto.Value{}
	})

	if err != nil {
		return err
	}

	err = ctx.ScriptVm.Set("napEnvSet", func(call otto.FunctionCall) otto.Value {
		ctx.EnvironmentVariables[call.Argument(0).String()] = call.Argument(1).String()

		return otto.Value{}
	})

	if err != nil {
		return err
	}

	err = ctx.ScriptVm.Set("napEnvGet", func(call otto.FunctionCall) otto.Value {
		result, _ := ctx.ScriptVm.ToValue(ctx.EnvironmentVariables[call.Argument(0).String()])
		return result
	})

	if err != nil {
		return err
	}

	err = ctx.ScriptVm.Set("napFail", func(call otto.FunctionCall) otto.Value {
		message := call.Argument(0).String()

		ctx.ScriptFailure = true
		ctx.ScriptFailureMessage = message

		return otto.Value{}
	})

	if err != nil {
		return err
	}

	_, err = ctx.ScriptVm.Run(`
var nap = { 
	env: { 
		get: napEnvGet, 
		set: napEnvSet
	}, 
	run: napRun,
	fail: napFail
};

napEnvGet = undefined;
napEnvSet = undefined;
napRun = undefined;
napFail = undefined;`)

	if err != nil {
		return err
	}

	return nil
}
