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

assert.go - this file contains types and logic for evaluating assertions
*/
package napassert

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Assert struct {
	Query       string
	Predicate   string
	Expectation string
}

func GetPredicates() []string {
	return []string{
		"==",
		">",
		">=",
		"<",
		"<=",
		"matches",
		"contains",
		"startswith",
		"endswith",
	}
}

func NewAssert(query string, predicate string, expectation string) *Assert {
	assert := new(Assert)

	assert.Query = query
	assert.Predicate = predicate
	assert.Expectation = expectation

	return assert
}

func AssertResponse(assertion *Assert, actual string) error {
	query := assertion.Query
	predicate := assertion.Predicate
	expectation := assertion.Expectation

	switch predicate {
	case "==":
		floatActual, err := strconv.ParseFloat(actual, 64)
		if err != nil && actual != expectation {
			return fmt.Errorf("%s of \"%s\" did not equal expected value \"%s\"", query, actual, expectation)
		} else {
			floatExpectation, err := strconv.ParseFloat(expectation, 64)
			if err != nil {
				return fmt.Errorf("%s of \"%s\" did not equal expected value \"%s\"", query, actual, expectation)
			} else if floatActual != floatExpectation {
				return fmt.Errorf("%s of \"%f\" did not equal expected value \"%f\"", query, floatActual, floatExpectation)
			}
		}
	case "<":
		floatActual, err := strconv.ParseFloat(actual, 64)
		if err != nil {
			return fmt.Errorf("%s of \"%s\" is not a number and cannot be compared using a less-than predicate", query, actual)
		}
		floatAssertValue, err := strconv.ParseFloat(expectation, 64)
		if err != nil {
			return fmt.Errorf("%s of \"%s\" is not a number and cannot be compared using a less-than predicate", query, actual)
		}

		if floatActual >= floatAssertValue {
			return fmt.Errorf("%s of \"%s\" is not less than expected value \"%s\"", query, actual, expectation)
		}
	case "<=":
		floatActual, err := strconv.ParseFloat(actual, 64)
		if err != nil {
			return fmt.Errorf("%s of \"%s\" is not a number and cannot be compared using a less-than predicate", query, actual)
		}
		floatAssertValue, err := strconv.ParseFloat(expectation, 64)
		if err != nil {
			return fmt.Errorf("%s of \"%s\" is not a number and cannot be compared using a less-than predicate", query, actual)
		}

		if floatActual > floatAssertValue {
			return fmt.Errorf("%s of \"%s\" is not less than expected value \"%s\"", query, actual, expectation)
		}
	case ">":
		floatActual, err := strconv.ParseFloat(actual, 64)
		if err != nil {
			return fmt.Errorf("%s of \"%s\" is not a number and cannot be compared using a less-than predicate", query, actual)
		}
		floatAssertValue, err := strconv.ParseFloat(expectation, 64)
		if err != nil {
			return fmt.Errorf("%s of \"%s\" is not a number and cannot be compared using a less-than predicate", query, actual)
		}

		if floatActual <= floatAssertValue {
			return fmt.Errorf("%s of \"%s\" is not less than expected value \"%s\"", query, actual, expectation)
		}
	case ">=":

		floatActual, err := strconv.ParseFloat(actual, 64)
		if err != nil {
			return fmt.Errorf("%s of \"%s\" is not a number and cannot be compared using a less-than predicate", query, actual)
		}
		floatAssertValue, err := strconv.ParseFloat(expectation, 64)
		if err != nil {
			return fmt.Errorf("%s of \"%s\" is not a number and cannot be compared using a less-than predicate", query, actual)
		}

		if floatActual < floatAssertValue {
			return fmt.Errorf("%s of \"%s\" is not less than expected value \"%s\"", query, actual, expectation)
		}
	case "matches":
		re := regexp.MustCompile(expectation)
		if !re.MatchString(actual) {
			return fmt.Errorf("%s of \"%s\" does not match expression /%s/", query, actual, expectation)
		}
	case "contains":
		if !strings.Contains(actual, expectation) {
			return fmt.Errorf("%s of \"%s\" does not contain string \"%s\"", query, actual, expectation)
		}
	case "startswith":
		if !strings.HasPrefix(actual, expectation) {
			return fmt.Errorf("%s of \"%s\" does not start with \"%s\"", query, actual, expectation)
		}
	case "endswith":
		if !strings.HasSuffix(actual, expectation) {
			return fmt.Errorf("%s of \"%s\" does not end with \"%s\"", query, actual, expectation)
		}
	default:
		return fmt.Errorf("Unrecognized predicate \"%s\"", predicate)
	}

	return nil
}
