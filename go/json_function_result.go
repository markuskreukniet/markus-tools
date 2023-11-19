package main

import "encoding/json"

type FunctionResult struct {
	Result       any
	ErrorMessage string
}

func createFunctionResultWithEmptyStringResult() FunctionResult {
	return FunctionResult{
		Result:       "",
		ErrorMessage: "",
	}
}

func jsonMarshalWithFallbackJSONError(nonJSON any) string {
	jsonBytes, err := json.Marshal(nonJSON)
	if err != nil {
		// This JSON string should match FunctionResult.
		// We can't use error.Error() for the message since that might turn the string into an invalid JSON.
		return `{"Result": null, "ErrorMessage": "json.Marshal error"}`
	}
	return string(jsonBytes)
}
