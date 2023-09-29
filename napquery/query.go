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

query.go - this file contains logic for evaluating queries
*/
package napquery

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PaesslerAG/jsonpath"

	"github.com/davesheldon/nap/napscript"
)

func Eval(query string, vmData *napscript.VmHttpData) (string, error) {
	if vmData == nil || vmData.Response == nil {
		// return empty here instead of erroring in case this assert is testing for absence of a value
		return "", nil
	}

	jsonExpression, isJsonPath := strings.CutPrefix(query, "jsonpath ")
	if isJsonPath {
		body := vmData.Response.JsonBody
		value, err := jsonpath.Get(jsonExpression, body)

		if err != nil {
			return "", err
		}

		return fmt.Sprint(value), nil
	}

	header, isHeader := strings.CutPrefix(query, "header ")
	if isHeader {
		if vmData.Response.Headers == nil {
			return "", nil
		}

		value := vmData.Response.Headers[header]

		return strings.Join(value, ","), nil
	}

	if query == "status" {
		return strconv.Itoa(vmData.Response.StatusCode), nil
	}

	if query == "duration" {
		return strconv.FormatInt(vmData.Response.ElapsedMs, 10), nil
	}

	if query == "body" {
		return vmData.Response.Body, nil
	}

	return "", fmt.Errorf("Query \"%s\" not recognized.", query)
}
