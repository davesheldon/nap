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

routine.go - this data structure represents a runnable set of instructions (a routine)
*/
package naproutine

import (
	"fmt"
	"github.com/davesheldon/nap/naprequest"
	"github.com/davesheldon/nap/napcontext"
	"github.com/kennygrant/sanitize"
	"gopkg.in/yaml.v2"
	"os"
	"path"
    "path/filepath"
	"strings"
	"sync"
	"time"
)

type Routine struct {
	Steps []*RoutineStep
}

type RoutineStep struct {
	Run string
}

func NewStep(run string) *RoutineStep {
	step := new(RoutineStep)
	step.Run = run
	return step
}

func NewRoutine(context *napcontext.Context, ...runs string) *Routine {
	routine := new(Routine)
	routine.Context = context

	for _, run := range runs {
		routine.steps = routine.steps.append(NewStep(run))
	}
}

func LoadFromPath(path string, ctx *napcontext.Context) (*Routine, error) {
	
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	dataAsString := string(data)

	for k, v := range ctx.EnvironmentVariables {
		variable := fmt.Sprintf("${%s}", k)
		dataAsString = strings.ReplaceAll(dataAsString, variable, v)
	}

	data = []byte(dataAsString)

	routine := parse(data)

	return routine
}

func parse(data []byte) (*Routine, error) {
	r := Routine{}
	err := yaml.Unmarshal(data, &r)

	if err != nil {
		return nil, err
	}

	return &r, nil
}
