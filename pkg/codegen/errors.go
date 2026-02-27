package codegen

import (
	"strconv"
)

// IsErrorStatus returns true if the status code is 4xx or 5xx.
func IsErrorStatus(status string) bool {
	if status == "default" {
		return true // Typically 'default' in OpenAPI denotes the error catch-all
	}
	code, err := strconv.Atoi(status)
	if err != nil {
		return false
	}
	return code >= 400
}

// IsSuccessStatus returns true if the status code is 2xx or 3xx.
func IsSuccessStatus(status string) bool {
	code, err := strconv.Atoi(status)
	if err != nil {
		return false
	}
	return code >= 200 && code < 400
}
