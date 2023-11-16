package main

import "encoding/json"

type FunctionResult struct {
	Result       interface{}
	ErrorMessage string
}

func jsonMarshalWithFallbackJSONError(nonJSON string) string {
	jsonBytes, err := json.Marshal(nonJSON)
	if err != nil {
		// This JSON string should match FunctionResult.
		// We can't use error.Error() for the message since that might turn the string into an invalid JSON.
		return `{"Result": null, "ErrorMessage": "json.Marshal error"}`
	}
	return string(jsonBytes)
}
