package main

import "encoding/json"

type FunctionResult struct {
	Result       interface{}
	ErrorMessage string
}

func jsonMarshalWithFallbackJSONError(nonJSON string) string {
	jsonBytes, err := json.Marshal(nonJSON)
	if err != nil {
		// TODO: comment
		return `{"Result": null, "ErrorMessage": "test"}`
	}
	return string(jsonBytes)
}
