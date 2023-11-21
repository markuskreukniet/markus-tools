package main

import "encoding/json"

type FunctionResult struct {
	Result       any
	ErrorMessage string
}

func createFunctionResult(result any, errorMessage string) FunctionResult {
	return FunctionResult{
		Result:       result,
		ErrorMessage: errorMessage,
	}
}

func defaultJSONFunctionResult() string {
	return jsonMarshalWithFallbackJSONError(createFunctionResult(nil, ""))
}

func resultToJSONFunctionResult(result any) string {
	return jsonMarshalWithFallbackJSONError(createFunctionResult(result, ""))
}

func errorMessageToJSONFunctionResult(errorMessage string) string {
	return jsonMarshalWithFallbackJSONError(createFunctionResult(nil, errorMessage))
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
