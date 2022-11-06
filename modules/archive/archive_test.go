package archive

import (
	"testing"

	"github.com/jarcoal/httpmock"

	"userstyles.world/models"
)

func TestIsFromArchive(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name, args string
		want       bool
	}{
		{"Test USw URL", "https://userstyles.world/api/style/1.user.css", false},
		{"Test bad URL", "https://github.com/userstyles-world/tweaks/raw/main/tweaks.user.styl", false},
		{"Test cdn URL", "https://cdn.jsdelivr.net/gh/33kk/uso-archive@flomaster/data/usercss/235518.user.css", true},
		{"Test hub URL", "https://raw.githubusercontent.com/33kk/uso-archive/flomaster/data/usercss/235518.user.css", true},
		{"Test org URL", "https://raw.githubusercontent.com/uso-archive/data/flomaster/data/usercss/235518.user.css", true},
		{"Test old URL", "https://uso-archive.surge.sh/?style=235518", true},
		{"Test new URL", "https://uso.kkx.one/style/235518", true},
	}

	for _, c := range cases {
		if got := IsFromArchive(c.args); got != c.want {
			t.Fatalf("%q: want: %t got: %t", c.name, c.want, got)
		}
	}
}

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

func TestExtractID(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name, input, expected string
	}{
		{"valid ID", ArchiveURL + "usercss/1.user.css", "1"},
		{"invalid ID", ArchiveURL + "usercss/-1.user.css", ""},
		{"bad input", ArchiveURL + "foo bar baz", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := extractID(c.input)
			if err != nil && err != ErrFailedToExtractID {
				t.Errorf("got: %s\n", err)
				t.Errorf("exp: %s\n", ErrFailedToExtractID)
			}
			if got != c.expected {
				t.Errorf("got: %s\n", got)
				t.Errorf("exp: %s\n", c.expected)
			}
		})
	}
}
