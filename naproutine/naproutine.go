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

naproutine.go - this data structure represents a runnable set of instructions (a routine)
*/
package naproutine

import (
	"fmt"
	"github.com/davesheldon/nap/naprequest"
	"github.com/kennygrant/sanitize"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Routine struct {
	Name  string
	Steps []*RoutineStep
}

type RoutineStep struct {
	Type                  string
	Name                  string
	ExpectStatusCode      string
	ExpectHeaders         map[string]string
	ExpectResponseContent string
	ExpectJson            string
}

func LoadByName(name string, environmentVariables map[string]string) (*Routine, error) {
	fileName := path.Join("routines", sanitize.BaseName(name)+".yml")

	data, err := os.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	dataAsString := string(data)

	for k, v := range environmentVariables {
		variable := fmt.Sprintf("${%s}", k)
		dataAsString = strings.ReplaceAll(dataAsString, variable, v)
	}

	data = []byte(dataAsString)

	return parse(data)
}

func parse(data []byte) (*Routine, error) {
	r := Routine{}
	err := yaml.Unmarshal(data, &r)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (routine *Routine) Run(environmentVariables map[string]string, parentStep *RoutineStep, ch chan *RoutineStepResult) *RoutineResult {
	result := new(RoutineResult)
	result.Name = routine.Name
	result.StartTime = time.Now()

	waitCount := 0

	childWg := new(sync.WaitGroup)
	childWg.Add(1)

	childCh := make(chan *RoutineStepResult)

	for _, step := range routine.Steps {
		var stepResult *RoutineStepResult

		if step.Type == "request" {
			request, err := naprequest.LoadByName(step.Name, environmentVariables)

			if err != nil {
				stepResult = StepError(step, err)
			} else {
				stepResult = StepRequestResult(step, request.Run())
			}
		}

		if step.Type == "routine" {
			subroutine, err := LoadByName(step.Name, environmentVariables)

			if err != nil {
				stepResult = StepError(step, err)
			} else {

				waitCount = waitCount + 1

				go subroutine.Run(environmentVariables, step, childCh)
				continue
			}
		}

		result.StepResults = append(result.StepResults, stepResult)
	}

	for i := 0; i < waitCount; i++ {
		deferredResult := <-childCh
		result.StepResults = append(result.StepResults, deferredResult)
	}

	if ch != nil {
		ch <- StepSubroutineResult(parentStep, result)
	}

	result.EndTime = time.Now()

	return result
}
