package utils

import "reflect"

// Checks if every entry of slice fulfills condition.
func Every(arr interface{}, cond func(interface{}) bool) bool {
	contentValue := reflect.ValueOf(arr)

	for i := 0; i < contentValue.Len(); i++ {
		if content := contentValue.Index(i); !cond(content.Interface()) {
			return false
		}
	}
	return true
}

// Check if array contains certain entry.
func Contains(arr []string, entry string) bool {
	for _, possibleEntry := range arr {
		if possibleEntry == entry {
			return true
		}
	}
	return false
}
