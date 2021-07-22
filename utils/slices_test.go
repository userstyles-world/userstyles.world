package utils

import "testing"

// sliceEqual is a helper function to check if the 2 slices are equal.
func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Test `Filter` to only keep the elements that satisfy the predicate.
func TestFilter(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		input     []int
		predicate func(interface{}) bool
		expected  []int
	}{
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(i interface{}) bool {
			return i.(int)%2 == 0
		}, []int{2, 4, 6, 8, 10}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(i interface{}) bool {
			return i.(int)%2 != 0
		}, []int{1, 3, 5, 7, 9}},
	}

	for _, tt := range tests {
		actual := Filter(tt.input, tt.predicate)
		if !sliceEqual(actual.([]int), tt.expected) {
			t.Errorf("expected %v, got %v", tt.expected, actual)
		}
	}
}

// Test `ContainsString` to return true if the given []string slice contains the given string.
func TestContainsString(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		input    []string
		entry    string
		expected bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "c", true},
		{[]string{"a", "b", "c"}, "d", false},
	}

	for _, tt := range tests {
		actual := ContainsString(tt.input, tt.entry)
		if actual != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, actual)
		}
	}
}

func containInString(input []string) func(name string) bool {
	return func(name string) bool {
		return ContainsString(input, name)
	}
}

// Test `EveryString` to return true if the given []string pass the condition.
func TestEveryString(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		input     []string
		predicate func(string) bool
		expected  bool
	}{
		{[]string{"a", "b", "c"}, containInString([]string{"a", "b", "c"}), true},
		{[]string{"a", "b", "c"}, containInString([]string{"a", "b", "d"}), false},
		{[]string{"a", "b", "c"}, containInString([]string{"a", "b", "d"}), false},
		{[]string{"a", "b", "c"}, containInString([]string{"a", "b", "c"}), true},
	}

	for _, tt := range tests {
		actual := EveryString(tt.input, tt.predicate)
		if actual != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, actual)
		}
	}
}
