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

request.go - this data structure represents a runnable HTTP request
*/
package naprequest

import (
	"bytes"
	"fmt"
	"github.com/kennygrant/sanitize"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Request struct {
	Path              string
	Verb              string
	Body              string
	Headers           map[string]string
	PreRequestScript  string
	PostRequestScript string
	PreRequestScriptFile  string
	PostRequestScriptFile string
}

func parse(data []byte) (*Request, error) {
	r := Request{}
	err := yaml.Unmarshal(data, &r)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func LoadFromPath(path string, ctx *napcontext.Context) (*Request, error) {
	fileName := path.Join("requests", sanitize.BaseName(name)+".yml")

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

func (r *Request) PrintStats() {
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