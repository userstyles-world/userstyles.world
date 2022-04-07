package utils

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"userstyles.world/models"
	"userstyles.world/modules/errors"

	libError "errors"
)

func TestArchiveImporting(t *testing.T) {
	t.Parallel()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "184563"
	httpmock.RegisterResponder("GET", DataURL+id+".json",
		httpmock.NewStringResponder(200, `{"info": {"description": "finally something new.", "additionalInfo": null, "category": "roblox"}, "screenshots": {"main": { "name": "184563_after.png"}}}`))

	httpmock.RegisterResponder("GET", StyleURL+id+".user.css",
		httpmock.NewStringResponder(200, `@-moz-document url-prefix(\"https://www.roblox.com/\") { * { display: none !important; } }`))

	data, err := ImportFromArchive(StyleURL+id+".user.css", models.APIUser{
		ID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if data.Description != "finally something new." {
		t.Fatal("Wrong description")
	}
	if data.Category != "roblox" {
		t.Fatal("Wrong category")
	}
	if data.Preview != PreviewURL+"184563_after.png" {
		t.Fatal("Wrong preview URL")
	}
	if data.Archived {
		t.Fatal("Should not be archived")
	}
	if data.Original != StyleURL+id+".user.css" {
		t.Fatal("Wrong original URL")
	}
}

func TestExtractingID(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc     string
		a        string
		expected any
	}{
		{"TestCorrectURL", StyleURL + "123.user.css", "123"},
		{"TestMaybeCorrectURL", StyleURL + "-123.user.css", "-123"},
		{"TestIncorrectURL", "What_Even-Is  This?!", errors.ErrStyleNotFromUSO},
	}

	for _, c := range cases {
		actual, err := extractID(c.a)
		if expecErr, ok := c.expected.(error); ok {
			if !libError.Is(err, expecErr) {
				t.Errorf("%s: Expected error %v, got %v", c.desc, expecErr, err)
			}
		} else if actual != c.expected {
			t.Fatalf("%s: expected: %s got: %s",
				c.desc, c.expected, actual)
		}
	}
}
