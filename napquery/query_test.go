package napquery_test

import (
	"testing"

	"github.com/davesheldon/nap/napquery"
	"github.com/davesheldon/nap/napscript"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func TestQueries(t *testing.T) {

	data := mockVmHttpData()

	tests := map[string]struct {
		query       string
		expectation []interface{}
	}{
		"body - default": {
			query:       "body",
			expectation: []any{data.Response.Body},
		},
		"jsonpath - simple": {
			query:       "jsonpath $.results[0].name",
			expectation: []any{"One"},
		},
		"jsonpath - length": {
			query:       "jsonpath $.results.length()",
			expectation: []any{6},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := napquery.Eval(test.query, data)
			if err != nil {
				t.Errorf("%T: %e", err, err)
			} else if actual[0] != test.expectation[0] {
				t.Errorf("Expected %v, got %v", test.expectation, actual)
			}
		})
	}
}

func mockVmHttpData() *napscript.VmHttpData {
	data := new(napscript.VmHttpData)
	data.Response = new(napscript.VmHttpResponse)
	data.Response.Headers = make(map[string][]any)
	data.Response.Headers["Content-Type"] = []any{"application/json"}
	data.Response.Body = `{ "results": [
		{
			"name": "One",
			"value": 1,
			"type": "Number"
		},
		{
			"name": "Two",
			"value": 2,
			"type": "Number"
		},
		{
			"name": "Three",
			"value": 3,
			"type": "Number"
		},
		{
			"name": "Red",
			"value": "#FF0000",
			"type": "Color"
		},
		{
			"name": "Blue",
			"value": "#0FF000",
			"type": "Color"
		},
		{
			"name": "Green",
			"value": "#0000FF",
			"type": "Color"
		}
	] }`
	json.Unmarshal([]byte(data.Response.Body), &data.Response.JsonBody)
	return data
}
