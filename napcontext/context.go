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
	ScriptVm             *otto.Otto
	EnvironmentName      string
	EnvironmentVariables map[string]string
	ScriptFailure        bool
	ScriptFailureMessage string
	ScriptOutput         []string
	WorkingDirectory     string
}

func New(workingDirectory string, environmentName string, environmentVariables map[string]string) *Context {
	ctx := new(Context)

	ctx.EnvironmentVariables = make(map[string]string)

	for k, v := range environmentVariables {
		ctx.EnvironmentVariables[k] = v
	}

	ctx.ScriptVm = otto.New()

	ctx.WorkingDirectory = workingDirectory

	return ctx
}

func (ctx *Context) Clone(workingDirectory string) *Context {
	return New(workingDirectory, ctx.EnvironmentName, ctx.EnvironmentVariables)
}
