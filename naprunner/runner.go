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
	"encoding/json"
	"errors"
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
	routine := naproutine.NewRoutine(ctx, path)
	result := runRoutine(ctx, routine, nil, nil)
	populateErrors(result)
	return result
}

func runScript(ctx *napcontext.Context, path string) *naproutine.ScriptResult {
	data, err := os.ReadFile(path)
	if err != nil {
		return naproutine.ScriptResultError(err)
	}

	script := string(data)

	return runScriptInline(ctx, script)
}

func runScriptInline(ctx *napcontext.Context, script string) *naproutine.ScriptResult {
	result := new(naproutine.ScriptResult)

	result.StartTime = time.Now()
	_, err := ctx.ScriptVm.Run(script)
	result.EndTime = time.Now()

	output := ctx.ScriptOutput
	ctx.ScriptOutput = []string{}

	result.Error = err
	result.ScriptOutput = output

	if result.Error == nil && ctx.ScriptFailure {
		result.Error = fmt.Errorf(ctx.ScriptFailureMessage)
	}

	return result
}

func runRequest(ctx *napcontext.Context, request *naprequest.Request) *naprequest.RequestResult {
	result := new(naprequest.RequestResult)
	result.Request = request

	if len(request.PreRequestScript) > 0 {
		scriptResult := runScriptInline(ctx, request.PreRequestScript)

		if scriptResult.Error != nil {
			result.Error = fmt.Errorf("Pre-Request Script Error: %w", scriptResult.Error)
			return result
		}
	}

	if len(request.PreRequestScriptFile) > 0 {
		scriptResult := runScript(ctx, request.PreRequestScriptFile)

		if scriptResult.Error != nil {
			result.Error = fmt.Errorf("Pre-Request Script File Error: %w", scriptResult.Error)
			return result
		}
	}

	result.StartTime = time.Now()

	result.HttpResponse, result.Error = executeHttp(request)

	if result.HttpResponse != nil && result.HttpResponse.StatusCode >= 400 {
		result.Error = errors.New(fmt.Sprintf("HTTP %s", result.HttpResponse.Status))
	}

	result.EndTime = time.Now()

	err := setVmHttpData(ctx, result)
	if err != nil {
		result.Error = fmt.Errorf("error setting js http result: %w", err)
		return result
	}

	if len(request.PostRequestScript) > 0 {
		scriptResult := runScriptInline(ctx, request.PostRequestScript)

		if scriptResult.Error != nil {
			result.Error = fmt.Errorf("Post-Request Script Error: %w", scriptResult.Error)
			return result
		}
	}

	if len(request.PostRequestScriptFile) > 0 {
		scriptResult := runScript(ctx, request.PostRequestScriptFile)

		if scriptResult.Error != nil {
			result.Error = fmt.Errorf("Post-Request Script File Error: %w", scriptResult.Error)
			return result
		}
	}

	return result
}

func executeHttp(r *naprequest.Request) (*http.Response, error) {
	client := &http.Client{}

	if r.TimeoutSeconds > 0 {
		client.Timeout = time.Duration(r.TimeoutSeconds) * time.Second
	}

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

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func runRoutine(ctx *napcontext.Context, routine *naproutine.Routine, parentStep *naproutine.RoutineStep, ch chan *naproutine.RoutineStepResult) *naproutine.RoutineResult {
	result := new(naproutine.RoutineResult)
	result.Routine = routine
	result.StartTime = time.Now()

	waitCount := 0

	childCh := make(chan *naproutine.RoutineStepResult)

	err := setupVm(ctx)
	if err != nil {
		result.EndTime = time.Now()
		result.Errors = append(result.Errors, err)
		return result
	}

	for _, step := range routine.Steps {
		var stepResult *naproutine.RoutineStepResult
		stepResult = nil

		stepPath := filepath.Join(ctx.WorkingDirectory, step.Run)

		if !fileExists(stepPath) {
			stepResult = naproutine.StepError(step, fmt.Errorf("file doesn't exist: %s", stepPath))
			result.StepResults = append(result.StepResults, stepResult)
			break
		}

		stepType, err := peekType(stepPath)

		if err != nil {
			stepResult = naproutine.StepError(step, err)
			result.StepResults = append(result.StepResults, stepResult)
			break
		}

		if stepType == "request" {
			request, err := naprequest.LoadFromPath(stepPath, ctx)

			if err != nil {
				stepResult = naproutine.StepError(step, err)
				result.StepResults = append(result.StepResults, stepResult)
				break
			} else {
				stepResult = naproutine.StepRequestResult(step, runRequest(ctx, request))
			}
		}

		if stepType == "script" {
			stepResult = naproutine.StepScriptResult(step, runScript(ctx, stepPath))
		}

		if stepType == "routine" {
			subroutineCtx := ctx.Clone(filepath.Dir(stepPath))
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

func peekType(path string) (string, error) {
	if filepath.Ext(path) == ".js" {
		return "script", nil
	}

	var yamlMap map[string]interface{}

	data, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("type of file unclear: %s (cannot read file: %s)", path, err.Error())
	}

	err = yaml.Unmarshal(data, &yamlMap)

	if err != nil {
		return "", fmt.Errorf("type of file unclear: %s (cannot unmarshal: %s)", path, err.Error())
	}

	if val, ok := yamlMap["type"]; ok {
		return val.(string), nil
	} else {
		return "", fmt.Errorf("type of file unclear: %s", path)
	}
}

func setupVm(ctx *napcontext.Context) error {

	err := ctx.ScriptVm.Set("napRun", func(call otto.FunctionCall) otto.Value {

		path := call.Argument(0).String()

		result := RunPath(ctx, path)

		v, _ := ctx.ScriptVm.ToValue(result)
		return v
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
		vals := []string{}

		for _, v := range call.ArgumentList {
			vals = append(vals, v.String())
		}

		ctx.ScriptFailure = true
		ctx.ScriptFailureMessage = strings.Join(vals, " ")

		return otto.Value{}
	})

	if err != nil {
		return err
	}

	ctx.ScriptVm.Set("__log__", func(call otto.FunctionCall) otto.Value {
		vals := []string{}

		for _, v := range call.ArgumentList {
			vals = append(vals, v.String())
		}

		ctx.ScriptOutput = append(ctx.ScriptOutput, strings.Join(vals, " "))

		return otto.Value{}
	})

	ctx.ScriptVm.Run("console.log = __log__;")

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
napFail = undefined;
`)

	if err != nil {
		return err
	}

	return nil
}

func setVmHttpData(ctx *napcontext.Context, result *naprequest.RequestResult) error {
	data, err := mapVmHttpData(result)

	if err != nil {
		return err
	}

	jsData, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	// have to do it this way for now because otto won't serialize it properly
	_, err = ctx.ScriptVm.Run(fmt.Sprintf("nap.http = %s;", string(jsData)))

	if err != nil {
		return err
	}

	return nil
}

type VmHttpData struct {
	Request  *VmHttpRequest  `json:"request"`
	Response *VmHttpResponse `json:"response"`
}

type VmHttpRequest struct {
	Url     string            `json:"url"`
	Verb    string            `json:"verb"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

type VmHttpResponse struct {
	StatusCode int                 `json:"statusCode"`
	Status     string              `json:"status"`
	Body       string              `json:"body"`
	JsonBody   interface{}         `json:"jsonBody"`
	Headers    map[string][]string `json:"headers"`
}

func mapVmHttpData(result *naprequest.RequestResult) (*VmHttpData, error) {
	data := new(VmHttpData)

	data.Request = new(VmHttpRequest)
	data.Request.Url = result.Request.Path
	data.Request.Verb = result.Request.Verb
	data.Request.Body = result.Request.Body
	data.Request.Headers = result.Request.Headers

	if result.HttpResponse != nil {
		data.Response = new(VmHttpResponse)
		data.Response.StatusCode = result.HttpResponse.StatusCode
		data.Response.Status = result.HttpResponse.Status

		bodyBytes, err := io.ReadAll(result.HttpResponse.Body)

		if err != nil {
			return nil, err
		}

		data.Response.Body = string(bodyBytes)

		defer result.HttpResponse.Body.Close()

		// TODO: support multiple header values per key
		data.Response.Headers = map[string][]string{}

		for k, v := range result.HttpResponse.Header {
			if len(v) > 0 {
				data.Response.Headers[k] = v

				if k == "Content-Type" && strings.Contains(v[0], "json") {
					err = json.Unmarshal(bodyBytes, &data.Response.JsonBody)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return data, nil
}
