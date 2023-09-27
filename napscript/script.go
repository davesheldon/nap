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

script.go - this file contains data structures and logic for setting up and executing against the javascript vm
*/

package napscript

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/naprequest"
	"github.com/davesheldon/nap/naproutine"
	"github.com/robertkrimen/otto"
)

func SetupVm(ctx *napcontext.Context, runPath func(*napcontext.Context, string) *naproutine.RoutineResult) error {

	err := ctx.ScriptVm.Set("napRun", func(call otto.FunctionCall) otto.Value {

		path := call.Argument(0).String()

		result := runPath(ctx, path)

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

		for _, v := range vals {
			ctx.ScriptOutput = append(ctx.ScriptOutput, v)
		}

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

func SetVmHttpData(ctx *napcontext.Context, result *naprequest.RequestResult) (*VmHttpData, error) {
	data, err := MapVmHttpData(result)

	if err != nil {
		return nil, err
	}

	jsData, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	// have to do it this way for now because otto won't serialize it properly
	_, err = ctx.ScriptVm.Run(fmt.Sprintf("nap.http = %s;", string(jsData)))

	if err != nil {
		return nil, err
	}

	return data, nil
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
	ElapsedMs  int64
}

func MapVmHttpData(result *naprequest.RequestResult) (*VmHttpData, error) {
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
		data.Response.ElapsedMs = result.GetElapsedMs()

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
