package napassert_test

import (
	"testing"

	"github.com/davesheldon/nap/napassert"
)

func TestAsserts(t *testing.T) {
	tests := map[string]struct {
		assert     *napassert.Assert
		actual     string
		shouldPass bool
	}{
		"equality passing": {
			assert:     napassert.NewAssert("", "==", "abc"),
			actual:     "abc",
			shouldPass: true,
		},
		"equality not passing": {
			assert:     napassert.NewAssert("", "==", "abc"),
			actual:     "123",
			shouldPass: false,
		},
		"not-equality passing": {
			assert:     napassert.NewAssert("", "not ==", "abc"),
			actual:     "123",
			shouldPass: true,
		},
		"not-equality not passing": {
			assert:     napassert.NewAssert("", "not ==", "abc"),
			actual:     "abc",
			shouldPass: false,
		},
		"inequality passing": {
			assert:     napassert.NewAssert("", "!=", "abc"),
			actual:     "123",
			shouldPass: true,
		},
		"inequality not passing": {
			assert:     napassert.NewAssert("", "!=", "abc"),
			actual:     "abc",
			shouldPass: false,
		},
		"not-inequality passing": {
			assert:     napassert.NewAssert("", "not !=", "abc"),
			actual:     "abc",
			shouldPass: true,
		},
		"not-inequality not passing": {
			assert:     napassert.NewAssert("", "not !=", "abc"),
			actual:     "123",
			shouldPass: false,
		},
		"gt passing": {
			assert:     napassert.NewAssert("", ">", "1"),
			actual:     "2",
			shouldPass: true,
		},
		"gt not passing": {
			assert:     napassert.NewAssert("", ">", "2"),
			actual:     "1",
			shouldPass: false,
		},
		"not-gt passing": {
			assert:     napassert.NewAssert("", "not >", "2"),
			actual:     "1",
			shouldPass: true,
		},
		"not-gt not passing": {
			assert:     napassert.NewAssert("", "not >", "1"),
			actual:     "2",
			shouldPass: false,
		},
		"lt passing": {
			assert:     napassert.NewAssert("", "<", "2"),
			actual:     "1",
			shouldPass: true,
		},
		"lt not passing": {
			assert:     napassert.NewAssert("", "<", "1"),
			actual:     "2",
			shouldPass: false,
		},
		"not-lt passing": {
			assert:     napassert.NewAssert("", "not <", "1"),
			actual:     "2",
			shouldPass: true,
		},
		"not-lt not passing": {
			assert:     napassert.NewAssert("", "not <", "2"),
			actual:     "1",
			shouldPass: false,
		},
		"gte passing": {
			assert:     napassert.NewAssert("", ">=", "1"),
			actual:     "2",
			shouldPass: true,
		},
		"gte not passing": {
			assert:     napassert.NewAssert("", ">=", "2"),
			actual:     "1",
			shouldPass: false,
		},
		"not-gte passing": {
			assert:     napassert.NewAssert("", "not >=", "2"),
			actual:     "1",
			shouldPass: true,
		},
		"not-gte not passing": {
			assert:     napassert.NewAssert("", "not >=", "1"),
			actual:     "2",
			shouldPass: false,
		},
		"lte passing": {
			assert:     napassert.NewAssert("", "<=", "2"),
			actual:     "1",
			shouldPass: true,
		},
		"lte not passing": {
			assert:     napassert.NewAssert("", "<=", "1"),
			actual:     "2",
			shouldPass: false,
		},
		"not-lte passing": {
			assert:     napassert.NewAssert("", "not <=", "1"),
			actual:     "2",
			shouldPass: true,
		},
		"not-lte not passing": {
			assert:     napassert.NewAssert("", "not <=", "2"),
			actual:     "1",
			shouldPass: false,
		},
		"matches passing": {
			assert:     napassert.NewAssert("", "matches", "^test.+$"),
			actual:     "testing123",
			shouldPass: true,
		},
		"matches not passing": {
			assert:     napassert.NewAssert("", "matches", "^test.+$"),
			actual:     "test",
			shouldPass: false,
		},
		"not matches passing": {
			assert:     napassert.NewAssert("", "not matches", "^test.+$"),
			actual:     "test",
			shouldPass: true,
		},
		"not matches not passing": {
			assert:     napassert.NewAssert("", "not matches", "^test.+$"),
			actual:     "testing123",
			shouldPass: false,
		},
		"contains passing": {
			assert:     napassert.NewAssert("", "contains", "bc12"),
			actual:     "abc123",
			shouldPass: true,
		},
		"contains not passing": {
			assert:     napassert.NewAssert("", "contains", "abcd"),
			actual:     "abc123",
			shouldPass: false,
		},
		"not contains passing": {
			assert:     napassert.NewAssert("", "not contains", "abcd"),
			actual:     "abc123",
			shouldPass: true,
		},
		"not contains not passing": {
			assert:     napassert.NewAssert("", "not contains", "bc12"),
			actual:     "abc123",
			shouldPass: false,
		},
		"startswith passing": {
			assert:     napassert.NewAssert("", "startswith", "abc"),
			actual:     "abc123",
			shouldPass: true,
		},
		"startswith not passing": {
			assert:     napassert.NewAssert("", "startswith", "def"),
			actual:     "abc123",
			shouldPass: false,
		},
		"not startswith passing": {
			assert:     napassert.NewAssert("", "not startswith", "def"),
			actual:     "abc123",
			shouldPass: true,
		},
		"not startswith not passing": {
			assert:     napassert.NewAssert("", "not startswith", "abc"),
			actual:     "abc123",
			shouldPass: false,
		},
		"endswith passing": {
			assert:     napassert.NewAssert("", "endswith", "123"),
			actual:     "abc123",
			shouldPass: true,
		},
		"endswith not passing": {
			assert:     napassert.NewAssert("", "endswith", "456"),
			actual:     "abc123",
			shouldPass: false,
		},
		"not endswith passing": {
			assert:     napassert.NewAssert("", "not endswith", "456"),
			actual:     "abc123",
			shouldPass: true,
		},
		"not endswith not passing": {
			assert:     napassert.NewAssert("", "not endswith", "123"),
			actual:     "abc123",
			shouldPass: false,
		},
		"in passing": {
			assert:     napassert.NewAssert("", "in", "[1,2,3]"),
			actual:     "1",
			shouldPass: true,
		},
		"in not passing": {
			assert:     napassert.NewAssert("", "in", "[1,2,3]"),
			actual:     "4",
			shouldPass: false,
		},
		"not in passing": {
			assert:     napassert.NewAssert("", "not in", "[1,2,3]"),
			actual:     "4",
			shouldPass: true,
		},
		"not in not passing": {
			assert:     napassert.NewAssert("", "not in", "[1,2,3]"),
			actual:     "1",
			shouldPass: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := napassert.Execute(test.assert, test.actual) == nil

			if result != test.shouldPass {
				t.Errorf("%s: expected passing=%v, got %v", name, test.shouldPass, result)
			}
		})
	}
}
