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

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/naprequest"
)

type RoutineResult struct {
	Routine     *Routine
	StepResults []*RoutineStepResult
	StartTime   time.Time
	EndTime     time.Time
	Errors      []error
}

func (r *RoutineResult) GetElapsedMs() int64 {
	return r.EndTime.Sub(r.StartTime).Milliseconds()
}

func (result *RoutineResult) IsPassing() bool {
	if len(result.Errors) > 0 {
		return false
	}

	for _, stepResult := range result.StepResults {
		if len(stepResult.Errors) > 0 {
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

func (result *RoutineResult) Print(prefix string, context *napcontext.Context) {
	if prefix == "" && result.Routine != nil && result.Routine.Name != "" {
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Routine: %s\n", result.Routine.Name)
		fmt.Println("------------------------------------------------------")
	}

	fmt.Printf("%sElapsed: %dms, IsPassing: %t\n", prefix, result.GetElapsedMs(), result.IsPassing())

	for i, s := range result.StepResults {
		s.print(i, prefix, context)
	}
}

func (result *RoutineResult) GetPassFailCounts() (passed int, failed int) {
	passed = 0
	failed = 0

	for _, v := range result.StepResults {
		if v.SubroutineResult != nil {
			subPassed, subFailed := v.SubroutineResult.GetPassFailCounts()
			passed += subPassed
			failed += subFailed
			continue
		}

		if v.RequestResult != nil {

			if v.RequestResult.Error != nil {
				failed += 1
			} else {
				passed += 1
			}
			continue
		}

		if v.ScriptResult != nil {
			if v.ScriptResult.Error != nil {
				failed += 1
			} else {
				passed += 1
			}
			continue
		}
	}

	return passed, failed
}

func (stepResult *RoutineStepResult) getName() string {
	if stepResult.RequestResult != nil {
		return fmt.Sprintf("%s (%s)", stepResult.RequestResult.Request.Name, stepResult.Step.Run)
	}
	if stepResult.SubroutineResult != nil {
		return fmt.Sprintf("%s (%s)", stepResult.SubroutineResult.Routine.Name, stepResult.Step.Run)
	}

	return stepResult.Step.Run
}

func (stepResult *RoutineStepResult) print(i int, prefix string, context *napcontext.Context) {
	fmt.Printf("%sRun %d: %s\n", prefix, i+1, stepResult.getName())

	for _, error := range stepResult.Errors {
		fmt.Printf("  [ERROR] %s\n", error.Error())
	}

	if stepResult.RequestResult != nil {
		if stepResult.RequestResult.Error != nil {
			fmt.Printf("%s  [ERROR] %s\n", prefix, stepResult.RequestResult.Error.Error())
		} else {
			fmt.Printf("%s  Status: %s\n", prefix, stepResult.RequestResult.HttpResponse.Status)
			fmt.Printf("%s  Elapsed: %dms\n", prefix, stepResult.RequestResult.GetElapsedMs())
		}
	}

	if stepResult.ScriptResult != nil {
		if stepResult.ScriptResult.Error != nil {
			fmt.Printf("%s  [ERROR] %s\n", prefix, stepResult.ScriptResult.Error.Error())
		} else {
			for _, v := range stepResult.ScriptResult.ScriptOutput {
				fmt.Printf("%s  Output: %s\n", prefix, v)
			}

			fmt.Printf("%s  Elapsed: %dms\n", prefix, stepResult.ScriptResult.GetElapsedMs())
		}
	}

	if stepResult.SubroutineResult != nil {
		stepResult.SubroutineResult.Print(prefix+"  ", context)
	}
}

func StepError(step *RoutineStep, err error) *RoutineStepResult {
	stepResult := new(RoutineStepResult)
	stepResult.Step = step
	stepResult.Errors = append(stepResult.Errors, err)

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
	Errors           []error
}

type ScriptResult struct {
	ScriptOutput []string
	StartTime    time.Time
	EndTime      time.Time
	Error        error
}

func (r *ScriptResult) GetElapsedMs() int64 {
	return r.EndTime.Sub(r.StartTime).Milliseconds()
}

func ScriptResultError(err error) *ScriptResult {
	result := new(ScriptResult)
	result.Error = err
	return result
}
