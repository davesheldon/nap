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

requestresult.go - this data structure represents the results of an executed Nap request
*/
package naprequest

import (
	"net/http"
	"time"
)

type RequestResult struct {
	Name              string
	HttpResponse      *http.Response
	PreRequestResult  string
	PostRequestResult string
	StartTime         time.Time
	EndTime           time.Time
	Error             error
}

func (r *RequestResult) GetElapsedMs() int64 {
	return r.EndTime.Sub(r.StartTime).Milliseconds()
}

func ResultError(err error) *RequestResult {
	result := new(RequestResult)
	result.Error = err
	return result
}
