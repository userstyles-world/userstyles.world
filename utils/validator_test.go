package utils

import "testing"

// Ensure that Initialization of Validator doesn't panic.
func TestValidatorInit(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("Panic: %v", err)
		}
	}()
	InitializeValidator()
}

// Test displayName validation. Which only should allow `^[a-zA-Z0-9-_ ]+$`
// But shouldn't allow any other weird characters and have a minuium of 5 chars and max 20.
func TestDisplayName(t *testing.T) {
	t.Parallel()
	InitializeValidator()

	type testDisplayName struct {
		DisplayName string `validate:"displayName,min=5,max=20"`
	}

	cases := []struct {
		description string
		actual      testDisplayName
		expected    bool
	}{
		{"valid", testDisplayName{DisplayName: "abcde"}, true},
		{"valid numerics", testDisplayName{DisplayName: "abcde123"}, true},
		{"valid spaces", testDisplayName{DisplayName: "abcde ef"}, true},
		{"too short", testDisplayName{DisplayName: "ab"}, false},
		{"too long", testDisplayName{DisplayName: "abcdefghijklmnopqrstuvwxyz"}, false},
		{"invalid characters", testDisplayName{DisplayName: "abcdefghijklmnopqrstuvwxyz#$"}, false},
	}

	for _, c := range cases {
		if err := Validate().StructPartial(c.actual, "DisplayName"); (err == nil) != c.expected {
			t.Errorf("%s: expected %t, got %s", c.description, c.expected, err)
		}
	}
}

// Test usernames validation. Which only should allow `^[a-zA-Z0-9_]+$`
// But shouldn't allow any other weird characters and have a minuium of 5 chars and max 20.
func TestUsername(t *testing.T) {
	t.Parallel()
	InitializeValidator()

	type testUsername struct {
		Username string `validate:"username,min=3,max=32"`
	}

	cases := []struct {
		description string
		actual      testUsername
		expected    bool
	}{
		{"valid", testUsername{Username: "abc"}, true},
		{"valid numerics", testUsername{Username: "123"}, true},
		{"valid characters", testUsername{Username: "o-0"}, true},
		{"invalid spaces", testUsername{Username: "abcde ef"}, false},
		{"too short", testUsername{Username: "ab"}, false},
		{"too long", testUsername{Username: "abcdefghijklmnopqrstuvwxyz1234567890"}, false},
		{"invalid characters", testUsername{Username: "tes_"}, false},
	}

	for _, c := range cases {
		if err := Validate().StructPartial(c.actual, "Username"); (err == nil) != c.expected {
			t.Errorf("%s: expected %t, got %s", c.description, c.expected, err)
		}
	}
}
