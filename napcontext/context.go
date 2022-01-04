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

context.go - this data represents a runnable context with its own set of environment variables and javascript virtual machine
*/
package napcontext

import (
	"github.com/robertkrimen/otto"
)

type Context struct {
	ScriptVm otto.Otto
	EnvironmentVariables map[string]string
	ScriptFailure bool
	ScriptFailureMessage string
}

func New(environmentVariables map[string]string) (*Context, error) {
	ctx =: new(Context)

	for k, v := range environmentVariables {
		ctx.EnvironmentVariables[k] = v
	}
	
	ctx.ScriptVm, err := setupVm(ctx)

	return ctx, nil
}

func (ctx *Context) Clone() *Context {
	return New(ctx.EnvironmentVariables)
}

func setupVm(ctx *Context) (*otto.Otto, error) {
	vm := otto.New()

	err := vm.Set("napRun", func(call otto.FunctionCall) otto.Value {
		
		path := call.Argument(0).String()
		result := Run(ctx, path)

		// todo: process result

		return otto.Value{}
	})

	if err != nil {
		return nil, err
	}

	err = vm.Set("napEnvSet", func(call otto.FunctionCall) otto.Value {
		ctx.EnvironmentVariables[call.Argument(0).String()] = call.Argument(1).String()

		return otto.Value{}
	})

	if err != nil {
		return nil, err
	}

	err = vm.Set("napEnvGet", func(call otto.FunctionCall) otto.Value {
		result, _ := vm.ToValue(ctx.EnvironmentVariables[call.Argument(0).String()])
		return result
	})

	if err != nil {
		return nil, err
	}

	err = vm.Set("napFail", func(call otto.FunctionCall) otto.Value {
		message := call.Argument(0).String()

		ctx.ScriptFailure = true
		ctx.ScriptFailureMessage = message

		return otto.Value{}
	})

	if err != nil {
		return nil, err
	}

	_, err = vm.Run(`
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
		return nil, err
	}

	return vm, nil
}
