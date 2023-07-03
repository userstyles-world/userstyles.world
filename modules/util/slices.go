package util

import "reflect"

// EveryString function checks if every entry of slice fulfills condition.
func EveryString(arr []string, cond func(string) bool) bool {
	for _, entry := range arr {
		if !cond(entry) {
			return false
		}
	}
	return true
}

// ContainsString func check if array contains certain entry.
func ContainsString(arr []string, entry string) bool {
	for _, possibleEntry := range arr {
		if possibleEntry == entry {
			return true
		}
	}
	return false
}

// Filter an slice while "preserving" the type with the reflect package.
func Filter(arr any, cond func(any) bool) any {
	contentType := reflect.TypeOf(arr)
	contentValue := reflect.ValueOf(arr)

	newContent := reflect.MakeSlice(contentType, 0, 0)
	for i := 0; i < contentValue.Len(); i++ {
		if content := contentValue.Index(i); cond(content.Interface()) {
			newContent = reflect.Append(newContent, content)
		}
	}
	return newContent.Interface()
}

func ContainsError(arr []error, entry error) bool {
	for _, possibleEntry := range arr {
		if possibleEntry == entry {
			return true
		}
	}
	return false
}
