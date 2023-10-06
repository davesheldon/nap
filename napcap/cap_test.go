package napcap_test

import (
	"fmt"
	"testing"

	"github.com/davesheldon/nap/napcap"
	"github.com/davesheldon/nap/napcontext"
	"github.com/davesheldon/nap/napscript"
)

func TestCaptures(t *testing.T) {
	tests := map[string]struct {
		variable  string
		ctx       *napcontext.Context
		queryFunc func(query string, vmData *napscript.VmHttpData) (string, error)
	}{
		"set new variable": {
			variable:  "test1",
			ctx:       napcontext.New("", nil, make(map[string]string), nil, true),
			queryFunc: mockQuery("value1", nil),
		},
		"overwrite new variable": {
			variable:  "test1",
			ctx:       napcontext.New("", nil, map[string]string{"test1": "value1"}, nil, true),
			queryFunc: mockQuery("value2", nil),
		},
		"error": {
			variable:  "test1",
			ctx:       napcontext.New("", nil, make(map[string]string), nil, true),
			queryFunc: mockQuery("", fmt.Errorf("mock error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			napcap.Query = test.queryFunc
			queryResult, queryError := test.queryFunc("", nil)
			err := napcap.CaptureQuery(test.variable, "", test.ctx, nil)

			if err == nil && queryError == nil && test.ctx.EnvironmentVariables[test.variable] != queryResult {
				t.Errorf("Expected %s=%s, got %s", test.variable, queryResult, test.ctx.EnvironmentVariables[test.variable])
			} else if queryError != nil && err == nil {
				t.Errorf("Expected error, got nil")
			} else if queryError == nil && err != nil {
				t.Errorf("Expected nil error, got %e", err)
			}
		})
	}
}

func mockQuery(mockResult string, mockError error) func(query string, vmData *napscript.VmHttpData) (string, error) {
	q := func(query string, vmData *napscript.VmHttpData) (string, error) {
		return mockResult, mockError
	}

	return q
}
