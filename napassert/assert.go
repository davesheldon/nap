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
	"encoding/json"
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
		"!=",
		">",
		">=",
		"<",
		"<=",
		"matches",
		"contains",
		"startswith",
		"endswith",
		"in",
	}
}

func NewAssert(query string, predicate string, expectation string) *Assert {
	assert := new(Assert)

	assert.Query = query
	assert.Predicate = predicate
	assert.Expectation = expectation

	return assert
}

func Execute(assert *Assert, actual []any) error {
	query := assert.Query
	predicate := assert.Predicate
	expectation := assert.Expectation

	// if the first character is a !, this will reverse the predicate evaluation
	basePredicate, isNot := strings.CutPrefix(predicate, "not ")

	// If there IS a ! then we want a false result to the predicate
	desiredResult := !isNot

	// init result as a failure
	result := !desiredResult

	var compareValue any = ""

	if len(actual) > 0 {
		compareValue = actual[0]
	}

	switch basePredicate {
	case "==":
		// special case until we implement json parsing of expectations
		if expectation == "[]" && len(actual) == 0 {
			result = true
			break
		}

		result = compareValue == expectation

		if !result {
			if compareValue == "" {
				compareValue = 0
			}
			floatActual, err := strconv.ParseFloat(fmt.Sprint(compareValue), 64)
			if err != nil {
				break
			}

			floatExpectation, err := strconv.ParseFloat(expectation, 64)
			if err != nil {
				break
			}

			result = floatActual == floatExpectation
		}
	case "!=":
		// special case until we implement json parsing of expectations
		if expectation == "[]" && len(actual) == 0 {
			result = false
			break
		}

		result = compareValue != expectation

		if result {
			if compareValue == "" {
				compareValue = 0
			}
			floatActual, err := strconv.ParseFloat(fmt.Sprint(compareValue), 64)
			if err != nil {
				break
			}

			floatExpectation, err := strconv.ParseFloat(expectation, 64)
			if err != nil {
				break
			}

			result = floatActual != floatExpectation
		}
	case "<":
		if compareValue == "" {
			compareValue = 0
		}
		floatActual, err := strconv.ParseFloat(fmt.Sprint(compareValue), 64)
		if err != nil {
			break
		}
		floatAssertValue, err := strconv.ParseFloat(expectation, 64)
		if err != nil {
			break
		}

		result = floatActual < floatAssertValue
	case "<=":
		if compareValue == "" {
			compareValue = 0
		}
		floatActual, err := strconv.ParseFloat(fmt.Sprint(compareValue), 64)
		if err != nil {
			break
		}
		floatAssertValue, err := strconv.ParseFloat(expectation, 64)
		if err != nil {
			break
		}

		result = floatActual <= floatAssertValue
	case ">":
		if compareValue == "" {
			compareValue = 0
		}
		floatActual, err := strconv.ParseFloat(fmt.Sprint(compareValue), 64)
		if err != nil {
			break
		}
		floatAssertValue, err := strconv.ParseFloat(expectation, 64)
		if err != nil {
			break
		}

		result = floatActual > floatAssertValue
	case ">=":
		if compareValue == "" {
			compareValue = 0
		}
		floatActual, err := strconv.ParseFloat(fmt.Sprint(compareValue), 64)
		if err != nil {
			break
		}
		floatAssertValue, err := strconv.ParseFloat(expectation, 64)
		if err != nil {
			break
		}

		result = floatActual >= floatAssertValue
	case "matches":
		re := regexp.MustCompile(expectation)
		result = re.MatchString(fmt.Sprint(compareValue))
	case "contains":
		result = strings.Contains(fmt.Sprint(compareValue), expectation)
	case "startswith":
		result = strings.HasPrefix(fmt.Sprint(compareValue), expectation)
	case "endswith":
		result = strings.HasSuffix(fmt.Sprint(compareValue), expectation)
	case "in":
		validValues := []interface{}{}
		data := []byte(expectation)
		err := json.Unmarshal(data, &validValues)
		if err != nil {
			break
		}

		in := false

		for _, val := range validValues {
			strVal := fmt.Sprint(val)
			in = strVal == compareValue

			if !in {
				if compareValue == "" {
					compareValue = 0
				}

				// string didn't compare, let's parse to float and try again
				floatVal, err := strconv.ParseFloat(strVal, 64)
				floatActual, err2 := strconv.ParseFloat(fmt.Sprint(compareValue), 64)

				if err == nil && err2 == nil {
					in = floatVal == floatActual
				}

			}

			if in {
				break
			}
		}
		result = in
	default:
		return fmt.Errorf("Unrecognized predicate \"%s\"", predicate)
	}

	if result != desiredResult {
		return fmt.Errorf("Assert failed \"%s => %s %s %s\"", query, actual, predicate, expectation)
	}

	return nil
}
