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

context.go - this data represents a runnable context with its own set of environment variables and javascript context
*/
package napcontext

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/davesheldon/nap/naputil"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Context struct {
	Environments         []string
	EnvironmentVariables map[string]string
	WorkingDirectory     string
	ScriptContext        *ScriptContext
	Cookies              []*http.Cookie

	progress  *mpb.Progress
	waitGroup *sync.WaitGroup
	quiet     bool
}

func New(workingDirectory string, environments []string, environmentVariables map[string]string, wg *sync.WaitGroup, quiet bool) *Context {
	ctx := new(Context)

	ctx.WorkingDirectory = workingDirectory
	ctx.Environments = environments
	ctx.EnvironmentVariables = map[string]string{}
	ctx.Cookies = []*http.Cookie{}

	for k, v := range environmentVariables {
		ctx.EnvironmentVariables[k] = v
	}

	ctx.waitGroup = wg
	ctx.quiet = quiet
	if !quiet {
		ctx.progress = mpb.New(mpb.WithWaitGroup(wg))
	}

	ctx.ScriptContext = newScriptContext()

	return ctx
}

func (old *Context) Clone(workingDirectory string) *Context {
	ctx := new(Context)

	ctx.Environments = old.Environments
	ctx.EnvironmentVariables = naputil.CloneMap(old.EnvironmentVariables)
	ctx.ScriptContext = newScriptContext()
	ctx.WorkingDirectory = workingDirectory
	ctx.progress = old.progress
	ctx.waitGroup = old.waitGroup
	ctx.quiet = old.quiet
	ctx.Cookies = append([]*http.Cookie{}, old.Cookies...)

	return ctx
}

type Progress struct {
	name  string
	steps int64
	bar   *mpb.Bar
}

func (ctx *Context) ProgressStart(name string, steps int64) *Progress {
	if ctx.quiet {
		return nil
	}
	ctx.waitGroup.Add(1)

	progress := new(Progress)
	progress.name = name
	progress.steps = steps
	progress.bar = ctx.progress.AddBar(steps,
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.CountersNoUnit("%d / %d", decor.WCSyncWidth), fmt.Sprintf("done (%d / %d).", steps, steps)),
		),
	)

	return progress
}

func (ctx *Context) ProgressIncrement(progress *Progress) {
	if ctx.quiet {
		return
	}
	progress.bar.Increment()

	if progress.bar.Current() == progress.steps {
		defer ctx.waitGroup.Done()
	}
}

func (ctx *Context) Complete() {
	if ctx.quiet {
		return
	}
	ctx.progress.Wait()
}
