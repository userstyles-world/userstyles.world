package utils

import "reflect"

// Every function checks if every entry of slice fulfills condition.
func Every(arr interface{}, cond func(interface{}) bool) bool {
	contentValue := reflect.ValueOf(arr)

	for i := 0; i < contentValue.Len(); i++ {
		if content := contentValue.Index(i); !cond(content.Interface()) {
			return false
		}
	}
	return true
}

// Contains func check if array contains certain entry.
func Contains(arr []string, entry string) bool {
	for _, possibleEntry := range arr {
		if possibleEntry == entry {
			return true
		}
	}
	return false
}

// Filter an slice while "preserving" the type with the reflect package.
func Filter(arr interface{}, cond func(interface{}) bool) interface{} {
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
