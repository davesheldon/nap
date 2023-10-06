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
	TimeoutSeconds        int `yaml:"timeoutSeconds"`
	Headers               map[string]string
	Body                  interface{}
	GraphQL               *GraphQLOptions `yaml:"graphql"`
	PreRequestScript      string          `yaml:"preRequestScript"`
	PostRequestScript     string          `yaml:"postRequestScript"`
	PreRequestScriptFile  string          `yaml:"preRequestScriptFile"`
	PostRequestScriptFile string          `yaml:"postRequestScriptFile"`
	Captures              map[string]string
	Asserts               []string
	Verbose               bool

	// aliases
	Url    string
	Method string
}

type GraphQLOptions struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
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
	request, err := parse(data)

	// check aliases
	if request != nil && len(request.Path) == 0 && len(request.Url) > 0 {
		request.Path = request.Url
	}

	if request != nil && len(request.Verb) == 0 && len(request.Method) > 0 {
		request.Verb = request.Method
	}

	return request, err
}

var expr = fmt.Sprintf("^(.+) (%s) \"?(.+)\"?$", strings.Join(napassert.GetPredicates(), "|"))
var re = regexp.MustCompile(expr)

func (request *Request) GetAsserts() ([]*napassert.Assert, error) {
	var asserts []*napassert.Assert = make([]*napassert.Assert, 0)
	for _, v := range request.Asserts {
		matches := re.FindStringSubmatch(v)
		if len(matches) < 4 {
			return nil, fmt.Errorf("Could not parse assert: %s", v)
		}

		query := matches[1]
		predicate := matches[2]
		expectation := matches[3]

		newQuery, queryHadNot := strings.CutSuffix(query, " not")
		if queryHadNot {
			query = newQuery
			predicate = "not " + predicate
		}

		assert := napassert.NewAssert(query, predicate, expectation)
		asserts = append(asserts, assert)
	}

	return asserts, nil
}
