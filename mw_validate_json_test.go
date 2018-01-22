package main

import (
	"encoding/json"
	"testing"

	"github.com/TykTechnologies/tyk/apidef"
	"github.com/TykTechnologies/tyk/test"
)

var testJsonSchema = `{
    "title": "Person",
    "type": "object",
    "properties": {
        "firstName": {
            "type": "string"
        },
        "lastName": {
            "type": "string"
        },
        "age": {
            "description": "Age in years",
            "type": "integer",
            "minimum": 0
        }
    },
    "required": ["firstName", "lastName"]
}`

func TestValidateJSONSchema(t *testing.T) {
	ts := newTykTestServer()
	defer ts.Close()

	buildAndLoadAPI(func(spec *APISpec) {
		updateAPIVersion(spec, "v1", func(v *apidef.VersionInfo) {
			json.Unmarshal([]byte(`[
				{
					"path": "/v",
					"method": "POST",
					"validate_with": `+testJsonSchema+`
				}
			]`), &v.ExtendedPaths.ValidateJSON)
		})

		spec.Proxy.ListenPath = "/"
	})

	ts.Run(t, []test.TestCase{
		{Method: "POST", Path: "/without_validation", Data: "{not_valid}", Code: 200},
		{Method: "POST", Path: "/v", Data: `{"age":23}`, BodyMatch: `firstName: firstName is required, lastName: lastName is required`, Code: 422},
		{Method: "POST", Path: "/v", Data: `[]`, BodyMatch: `Expected: object, given: array`, Code: 422},
		{Method: "POST", Path: "/v", Data: `not_json`, Code: 400},
		{Method: "POST", Path: "/v", Data: `{"age":23, "firstName": "Harry", "lastName": "Potter"}`, Code: 200},
	}...)
}
