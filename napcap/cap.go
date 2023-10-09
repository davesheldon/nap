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

capture.go - this file contains logic for evaluating captures
*/
package napcap

import (
	"fmt"

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/napquery"
	"github.com/davesheldon/nap/napscript"
)

var (
	Query = napquery.Eval
)

func CaptureQuery(variable string, query string, ctx *napcontext.Context, vmData *napscript.VmHttpData) error {
	actual, err := Query(query, vmData)

	if err != nil {
		return err
	}

	if len(actual) > 0 {
		ctx.EnvironmentVariables[variable] = fmt.Sprint(actual[0])
	}

	// todo: deal with multiple return values

	return nil
}
