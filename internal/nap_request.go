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
*/
package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type NapRequest struct {
	Name              string
	Path              string
	Verb              string
	Type              string
	Body              string
	Headers           map[string]string
	PreRequestScript  string
	PostRequestScript string
}

func ParseNapRequest(data []byte) (*NapRequest, error) {
	r := NapRequest{}
	err := yaml.Unmarshal(data, &r)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *NapRequest) GetResult() *NapRequestResult {
	result := new(NapRequestResult)

	result.StartTime = time.Now()
	result.HttpResponse, result.Error = r.executeHttp()
	result.EndTime = time.Now()

	return result
}

func (r *NapRequest) PrintStats() {
	fmt.Printf("\n\nRunning: %s\n", r.Name)
	fmt.Printf("Path: %s\n", r.Path)
	fmt.Printf("Verb: %s\n", r.Verb)

	for k, v := range r.Headers {
		fmt.Printf("(Header) %s: %s\n", k, v)
	}

	if len(r.Body) > 0 {
		fmt.Printf("Request Body: %s\n", r.Body)
	}
}

func (r *NapRequest) executeHttp() (*http.Response, error) {
	client := &http.Client{}

	var content io.Reader

	if len(r.Body) > 0 {
		content = bytes.NewBuffer([]byte(r.Body))
	} else {
		content = strings.NewReader("")
	}

	request, err := http.NewRequest(r.Verb, r.Path, content)

	if err != nil {
		return nil, err
	}

	for k, v := range r.Headers {
		request.Header.Add(k, v)
	}

	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
