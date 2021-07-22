package utils

import (
	"bytes"
	"testing"
	"time"

	"github.com/ohler55/ojg/oj"
	"userstyles.world/search"
)

func TestUnsafeString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc     string
		a        []byte
		expected string
	}{
		{"TestProperName", []byte("UserStyle Name"), "UserStyle Name"},
		{"TestMoreCharacters", []byte("What_Even-Is  This?!"), "What_Even-Is  This?!"},
		{"TestExtraCharacters", []byte("(Dark) Something [v1.2.3]"), "(Dark) Something [v1.2.3]"},
		{"TextExtraOfEverything", []byte(" Please---Get___Some   HELP!!! "), " Please---Get___Some   HELP!!! "},
		{"TestTypographicSymbols", []byte("暗い空 Dark Mode)"), "暗い空 Dark Mode)"},
	}

	for _, c := range cases {
		actual := UnsafeString(c.a)
		if actual != c.expected {
			t.Fatalf("%s: expected: %s got: %s",
				c.desc, c.expected, actual)
		}
	}
}

func TestUnsafeBytes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc     string
		expected []byte
		a        string
	}{
		{"TestProperName", []byte("UserStyle Name"), "UserStyle Name"},
		{"TestMoreCharacters", []byte("What_Even-Is  This?!"), "What_Even-Is  This?!"},
		{"TestExtraCharacters", []byte("(Dark) Something [v1.2.3]"), "(Dark) Something [v1.2.3]"},
		{"TextExtraOfEverything", []byte(" Please---Get___Some   HELP!!! "), " Please---Get___Some   HELP!!! "},
		{"TestTypographicSymbols", []byte("暗い空 Dark Mode)"), "暗い空 Dark Mode)"},
	}

	for _, c := range cases {
		actual := UnsafeBytes(c.a)
		if !bytes.Equal(actual, c.expected) {
			t.Fatalf("%s: expected: %s got: %s",
				c.desc, c.expected, actual)
		}
	}
}

type testStruct struct {
	Name string
}

func TestJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc     string
		a        interface{}
		expected []byte
	}{
		{"SimpleTest", testStruct{
			Name: "abcv",
		}, []byte(`{"name":"abcv"}`)},
		{"TestForMinimalStyle", search.MinimalStyle{
			Name:        "abcv",
			ID:          123,
			CreatedAt:   time.Unix(0, 0),
			UpdatedAt:   time.Unix(0, 0),
			Username:    "admin",
			DisplayName: "Admin",
			Description: "This is a description",
			Preview:     "https://example.com/preview.png",
			Notes:       "This is a note",
			Views:       99,
			Installs:    69,
		}, []byte(`{"created_at":0,"description":"This is a description","display_name":"Admin","id":123,"installs":69,"name":"abcv","notes":"This is a note","preview":"https://example.com/preview.png","updated_at":0,"username":"admin","views":99}`)},
	}

	for _, c := range cases {
		actual := UnsafeBytes(oj.JSON(c.a, &oj.Options{OmitNil: true, UseTags: true, Sort: true}))
		if !bytes.Equal(actual, c.expected) {
			t.Fatalf("%s: expected: %s got: %s",
				c.desc, c.expected, actual)
		}
	}
}

// Test the `EncodeToString` function
func TestBase64Encoding(t *testing.T) {
	t.Parallel()

	cases := []struct {
		a        string
		expected string
	}{
		// RFC 3548 examples
		{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l-"},
		{"\x14\xfb\x9c\x03\xd9", "FPucA9k"},
		{"\x14\xfb\x9c\x03", "FPucAw"},

		// RFC 4648 examples
		{"", ""},
		{"f", "Zg"},
		{"fo", "Zm8"},
		{"foo", "Zm9v"},
		{"foob", "Zm9vYg"},
		{"fooba", "Zm9vYmE"},
		{"foobar", "Zm9vYmFy"},

		// Wikipedia examples
		{"sure.", "c3VyZS4"},
		{"sure", "c3VyZQ"},
		{"sur", "c3Vy"},
		{"su", "c3U"},
		{"leasure.", "bGVhc3VyZS4"},
		{"easure.", "ZWFzdXJlLg"},
		{"asure.", "YXN1cmUu"},
		{"sure.", "c3VyZS4"},
	}

	for _, c := range cases {
		actual := EncodeToString([]byte(c.a))
		if actual != c.expected {
			t.Fatalf("expected: %s got: %s",
				c.expected, actual)
		}
	}
}
