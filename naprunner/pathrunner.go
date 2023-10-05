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
	"path"

	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/naproutine"
)

func RunPath(ctx *napcontext.Context, runPath string) *naproutine.RoutineResult {
	routine := naproutine.NewRoutine(ctx, "Job: "+path.Base(runPath), runPath)
	result := runRoutine(ctx, routine, nil, nil)
	populateErrors(result)

	return result
}
