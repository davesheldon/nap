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
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/davesheldon/nap/napassert"
	"github.com/davesheldon/nap/napcap"
	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/napquery"
	"github.com/davesheldon/nap/naprequest"
	"github.com/davesheldon/nap/napscript"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func runRequest(ctx *napcontext.Context, runPath string, request *naprequest.Request) *naprequest.RequestResult {
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

	response, err := executeHttp(request, ctx, filepath.Dir(runPath))

	result.HttpResponse = response
	if err != nil {
		result.Error = fmt.Errorf("Request failed to execute: %w", err)
		return result
	}

	result.EndTime = time.Now()

	if err := napscript.SetupVm(ctx, RunPath); err != nil {
		result.Error = fmt.Errorf("Error setting up js vm: %w", err)
		return result
	}

	vmData, err := napscript.SetVmHttpData(ctx, result)
	if err != nil {
		result.Error = fmt.Errorf("Error setting js http result: %w", err)
		return result
	}

	for variable, query := range request.Captures {
		err := napcap.CaptureQuery(variable, query, ctx, vmData)
		if err != nil {
			result.Error = err
			return result
		}
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

	asserts, err := request.GetAsserts(ctx)
	if err != nil {
		result.Error = err
		return result
	}

	for _, v := range asserts {

		actual, err := napquery.Eval(v.Query, vmData)

		if err != nil {
			result.Error = err
			return result
		}

		var testVal interface{} = nil

		if actual != nil && len(actual) > 0 {
			testVal = actual[0]
		}

		err = napassert.Execute(v, testVal)

		if err != nil {
			result.Error = err
			return result
		}
	}

	return result
}

func executeHttp(r *naprequest.Request, ctx *napcontext.Context, workingDirectory string) (*http.Response, error) {
	client := &http.Client{}

	if r.TimeoutSeconds > 0 {
		client.Timeout = time.Duration(r.TimeoutSeconds) * time.Second
	}

	var content io.Reader

	if r.GraphQL != nil {
		if strings.HasPrefix(r.GraphQL.Query, "@") {
			gqlPayloadFile := r.GraphQL.Query[1:]
			gqlFullPath := filepath.Join(workingDirectory, gqlPayloadFile)
			gqlData, err := os.ReadFile(gqlFullPath)
			if err != nil {
				return nil, err
			}
			r.GraphQL.Query = string(gqlData)
		}

		graphqlPayload, err := json.Marshal(&r.GraphQL)
		if err != nil {
			return nil, err
		}

		content = bytes.NewBuffer(graphqlPayload)
		r.Verb = "POST"
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		r.Headers["Content-Type"] = "application/json"
	} else if bodyAsString := fmt.Sprint(r.Body); r.Body != nil && len(bodyAsString) > 0 {
		if strings.HasPrefix(r.Headers["Content-Type"], "multipart/form-data") {
			bodyAsMap, ok := r.Body.(map[interface{}]interface{})

			if !ok {
				return nil, fmt.Errorf("Could not read form body.")
			}

			bodyAsStringMap := make(map[string]string)

			for k, v := range bodyAsMap {
				stringKey := fmt.Sprint(k)
				bodyAsStringMap[stringKey] = fmt.Sprint(v)
			}

			newHeader, formData, err := createFormData(bodyAsStringMap, workingDirectory)
			if err != nil {
				return nil, err
			}

			r.Headers["Content-Type"] = newHeader
			content = formData
		} else if strings.HasPrefix(bodyAsString, "@") {
			bodyFileName := bodyAsString[1:]
			pathToPayload := filepath.Join(workingDirectory, bodyFileName)
			file, err := os.ReadFile(pathToPayload)
			if err != nil {
				return nil, err
			}
			content = bytes.NewBuffer(file)
		} else {
			bodyAsString, ok := r.Body.(string)
			if ok {
				content = bytes.NewBuffer([]byte(bodyAsString))
			} else {
				// todo: would be nice to detect the content type and marshal to the right format here to convert e.g. YAML to JSON
				return nil, fmt.Errorf("Could not read body as string")
			}
		}
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

	if len(r.Cookies)+len(ctx.Cookies) > 0 {
		cookies := make([]*http.Cookie, 0)
		for k, v := range r.Cookies {
			cookies = append(cookies, &http.Cookie{Name: k, Value: v})
		}
		cookies = append(cookies, ctx.Cookies...)

		for _, v := range cookies {
			request.AddCookie(v)
		}
	}

	if r.Verbose {
		fmt.Println("REQUEST:")
		dump, err := httputil.DumpRequestOut(request, true)
		if err == nil {

			fmt.Println(string(dump))
		} else {
			fmt.Println(err)
		}
	}

	response, err := client.Do(request)
	ctx.Cookies = append(ctx.Cookies, response.Cookies()...)

	if r.Verbose {
		fmt.Println("RESPONSE:")
		dump, err := httputil.DumpResponse(response, true)
		if err == nil {
			fmt.Println(string(dump))
		} else {
			fmt.Println(err)
		}
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func createFormData(form map[string]string, workingDirectory string) (string, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			file, err := os.Open(path.Join(workingDirectory, val))
			if err != nil {
				return "", nil, err
			}
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil {
				return "", nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
		}
	}
	return mp.FormDataContentType(), body, nil
}
