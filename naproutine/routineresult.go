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

routineresult.go - this data structure represents a routine run result
*/
package naproutine

import (
	"fmt"
	"time"

	"github.com/davesheldon/nap/naprequest"
)

type RoutineResult struct {
	StepResults []*RoutineStepResult
	StartTime   time.Time
	EndTime     time.Time
	Error       error
}

func (r *RoutineResult) GetElapsedMs() int64 {
	return r.EndTime.Sub(r.StartTime).Milliseconds()
}

func (result *RoutineResult) IsPassing() bool {
	if result.Error != nil {
		return false
	}

	for _, stepResult := range result.StepResults {
		if stepResult.Error != nil {
			return false
		}

		if stepResult.RequestResult != nil && stepResult.RequestResult.Error != nil {
			return false
		}

		if stepResult.SubroutineResult != nil && !stepResult.SubroutineResult.IsPassing() {
			return false
		}
	}

	return true
}

func (result *RoutineResult) Print(prefix string) {
	fmt.Printf("%sElapsedMs: %d, IsPassing: %t\n", prefix, result.GetElapsedMs(), result.IsPassing())

	for i, s := range result.StepResults {
		s.print(i, prefix)
	}
}

func (stepResult *RoutineStepResult) print(i int, prefix string) {
	fmt.Printf("%sRun %d: %s\n", prefix, i+1, stepResult.Step.Run)
	if stepResult.Error != nil {
		fmt.Printf("- ERROR! %s", stepResult.Error.Error())
	}

	if stepResult.RequestResult != nil {
		if stepResult.RequestResult.Error != nil {
			fmt.Printf("%s- ERROR! %s\n", prefix, stepResult.RequestResult.Error.Error())
		} else {
			fmt.Printf("%s- status: %s\n", prefix, stepResult.RequestResult.HttpResponse.Status)
		}
	}

	if stepResult.SubroutineResult != nil {
		stepResult.SubroutineResult.Print(prefix + "- ")
	}
}

func StepError(step *RoutineStep, err error) *RoutineStepResult {
	stepResult := new(RoutineStepResult)
	stepResult.Step = step
	stepResult.Error = err

	return stepResult
}

func StepRequestResult(step *RoutineStep, requestResult *naprequest.RequestResult) *RoutineStepResult {
	stepResult := new(RoutineStepResult)
	stepResult.Step = step

	stepResult.RequestResult = requestResult

	return stepResult
}

func StepScriptResult(step *RoutineStep, scriptResult *ScriptResult) *RoutineStepResult {
	stepResult := new(RoutineStepResult)
	stepResult.Step = step
	stepResult.ScriptResult = scriptResult

	return stepResult
}

func StepSubroutineResult(step *RoutineStep, subroutineResult *RoutineResult) *RoutineStepResult {
	stepResult := new(RoutineStepResult)
	stepResult.Step = step

	stepResult.SubroutineResult = subroutineResult

	return stepResult
}

type RoutineStepResult struct {
	Step             *RoutineStep
	RequestResult    *naprequest.RequestResult
	SubroutineResult *RoutineResult
	ScriptResult     *ScriptResult
	Error            error
}

type ScriptResult struct {
	ScriptOutput []string
	Error        error
}

func ScriptResultError(err error) *ScriptResult {
	result := new(ScriptResult)
	result.Error = err
	return result
}
