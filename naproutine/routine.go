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
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/napenv"
	"github.com/davesheldon/nap/naputil"
	"gopkg.in/yaml.v2"
)

type Routine struct {
	Name  string
	Steps []*RoutineStep
}

type RoutineStep struct {
	Run        string
	Iterations interface{}
}

func (step *RoutineStep) GetIterations(ctx *napcontext.Context) ([]*napcontext.Context, error) {
	iterations := make([]*napcontext.Context, 0)

	var environmentsToLoad []string = make([]string, 0)

	switch step.Iterations.(type) {
	case string:
		if result, err := filepath.Glob(path.Join(ctx.WorkingDirectory, step.Iterations.(string))); err != nil {
			return nil, err
		} else {
			environmentsToLoad = result
		}
	case []interface{}:
		for _, v := range step.Iterations.([]interface{}) {
			switch v.(type) {
			case string:
				if result, err := filepath.Glob(path.Join(ctx.WorkingDirectory, v.(string))); err != nil {
					return nil, err
				} else {
					environmentsToLoad = append(environmentsToLoad, result...)
				}
			}
		}
	}

	if environmentsToLoad != nil {
		for _, v := range environmentsToLoad {
			iteration := ctx.Clone(ctx.WorkingDirectory)
			baseEnv := naputil.CloneMap(ctx.EnvironmentVariables)

			result, err := napenv.AddEnvironmentFromPath(ctx.WorkingDirectory, v, baseEnv)
			if err != nil {
				return nil, err
			}

			iteration.EnvironmentVariables = result

			iterations = append(iterations, iteration)
		}
	}

	if len(iterations) == 0 {
		iterations = append(iterations, ctx)
	}

	return iterations, nil
}

func NewStep(run string, iterations interface{}) *RoutineStep {
	step := new(RoutineStep)
	step.Run = run
	step.Iterations = iterations
	return step
}

func NewRoutine(ctx *napcontext.Context, name string, runs ...string) *Routine {
	routine := new(Routine)
	routine.Name = name

	for _, run := range runs {
		routine.Steps = append(routine.Steps, NewStep(run, nil))
	}

	return routine
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

	routine, err := parse(data)
	if err != nil {
		return nil, err
	}

	if routine.Name == "" {
		routine.Name = path
	}

	return routine, nil
}

func parse(data []byte) (*Routine, error) {
	r := Routine{}
	err := yaml.Unmarshal(data, &r)

	if err != nil {
		return nil, err
	}

	return &r, nil
}
