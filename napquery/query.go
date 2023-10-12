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

	"github.com/AsaiYusuke/jsonpath"
	"github.com/davesheldon/nap/napscript"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func evalJsonPath(expression string, data interface{}) ([]interface{}, error) {
	var config = jsonpath.Config{}
	config.SetAggregateFunction(`length`, func(params []interface{}) (interface{}, error) {
		if params == nil {
			return nil, nil
		}
		return len(params), nil
	})
	output, err := jsonpath.Retrieve(expression, data, config)

	// return empty for certain types of errors
	if err != nil {
		switch err.(type) {
		case jsonpath.ErrorMemberNotExist:
			return []interface{}{}, nil
		}
	}

	return output, err
}

func Eval(query string, vmData *napscript.VmHttpData) ([]any, error) {
	if vmData == nil || vmData.Response == nil {
		// return empty here instead of erroring in case this assert is testing for absence of a value
		return nil, nil
	}

	jsonExpression, isJsonPath := strings.CutPrefix(query, "jsonpath ")
	if isJsonPath {
		body := vmData.Response.JsonBody
		value, err := evalJsonPath(jsonExpression, body)

		if err != nil {
			return nil, err
		}

		return value, nil
	}

	header, isHeader := strings.CutPrefix(query, "header ")
	if isHeader {
		if vmData.Response.Headers == nil {
			return nil, nil
		}

		value := vmData.Response.Headers[header]

		return value, nil
	}

	cookie, isCookie := strings.CutPrefix(query, "cookie ")
	if isCookie {
		if strings.Contains(cookie, "[") {
			cookieParts := strings.Split(cookie, "[")

			targetCookie, ok := vmData.Response.Cookies[strings.TrimSpace(cookieParts[0])]
			if !ok {
				return []any{}, nil
			}

			cookieAttr := strings.Replace(strings.TrimSpace(cookieParts[1]), "]", "", 0)
			switch cookieAttr {
			case "Value":
				return []any{targetCookie.Value}, nil
			case "Expires":
				return []any{targetCookie.RawExpires}, nil
			case "Max-Age":
				return []any{targetCookie.MaxAge}, nil
			case "Domain":
				return []any{targetCookie.Domain}, nil
			case "Path":
				return []any{targetCookie.Path}, nil
			case "Secure":
				return []any{targetCookie.Secure}, nil
			case "HttpOnly":
				return []any{targetCookie.HttpOnly}, nil
			case "SameSite":
				return []any{targetCookie.SameSite}, nil
			}

		} else {
			targetCookie, ok := vmData.Response.Cookies[cookie]
			if !ok {
				return []any{}, nil
			}

			return []any{targetCookie.Value}, nil
		}
	}

	if query == "status" {
		return []any{strconv.Itoa(vmData.Response.StatusCode)}, nil
	}

	if query == "duration" {
		return []any{strconv.FormatInt(vmData.Response.ElapsedMs, 10)}, nil
	}

	if query == "body" {
		return []any{vmData.Response.Body}, nil
	}

	return nil, fmt.Errorf("Query \"%s\" not recognized.", query)
}
