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
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/davesheldon/nap/napassert"
	"github.com/davesheldon/nap/napcontext"
	"gopkg.in/yaml.v2"
)

type Request struct {
	Name                  string
	Path                  string
	Verb                  string
	Body                  string
	Headers               map[string]string
	PreRequestScript      string `yaml:"preRequestScript"`
	PostRequestScript     string `yaml:"postRequestScript"`
	PreRequestScriptFile  string `yaml:"preRequestScriptFile"`
	PostRequestScriptFile string `yaml:"postRequestScriptFile"`
	TimeoutSeconds        int    `yaml:"timeoutSeconds"`
	Captures              map[string]string
	Asserts               []string
	Verbose               bool
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
	return parse(data)
}

var expr = fmt.Sprintf("^(?P<Query>.+) (?P<Predicate>%s) \"?(?P<Value>.+)\"?$", strings.Join(napassert.GetPredicates(), "|"))
var re = regexp.MustCompile(expr)

func (request *Request) GetAsserts() []*napassert.Assert {
	var asserts []*napassert.Assert = make([]*napassert.Assert, 0)
	for _, v := range request.Asserts {
		matches := re.FindStringSubmatch(v)

		query := matches[1]
		predicate := matches[2]
		expectation := matches[3]

		assert := napassert.NewAssert(query, predicate, expectation)

		asserts = append(asserts, assert)
	}

	return asserts
}
