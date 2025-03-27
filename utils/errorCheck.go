package utils

import "strings"

func IsCriticalError(err error) (string, bool) {
	errors := []string{"openapi svc error", "500 Internal Server Error", "Please try again with a different amount or token pair."}
	for _, errValue := range errors {
		if strings.Contains(err.Error(), errValue) {
			return err.Error(), false
		}
	}
	return err.Error(), true
}
