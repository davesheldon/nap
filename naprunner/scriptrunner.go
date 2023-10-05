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
	"os"
	"time"

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/naproutine"
	"github.com/davesheldon/nap/napscript"
)

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

	if err := napscript.SetupVm(ctx, RunPath); err != nil {
		result.Error = err
		return result
	}

	result.StartTime = time.Now()
	_, err := ctx.ScriptContext.Vm.Run(script)
	result.EndTime = time.Now()

	output := ctx.ScriptContext.Output
	ctx.ScriptContext.Output = []string{}

	result.Error = err
	result.ScriptOutput = output

	if result.Error == nil && ctx.ScriptContext.Error != nil {
		result.Error = ctx.ScriptContext.Error
	}

	return result
}
