package models

import (
	"strings"
	"testing"
)

var reviewCases = []struct {
	name  string
	input *Review
	exp   error
}{
	{
		"both empty",
		NewReview(1, 1, "0", ""),
		errorReviewEmptyFields,
	},
	{
		"bad rating",
		NewReview(1, 1, "foo", "bar"),
		errorReviewOutOfRange,
	},
	{
		"bad comment",
		NewReview(1, 1, "5", strings.Repeat(";", 501)),
		errorReviewLongComment,
	},
	{
		"both set",
		NewReview(1, 1, "5", "foo"),
		nil,
	},
}

func TestReview_Validate(t *testing.T) {
	t.Parallel()
	for _, c := range reviewCases {
		t.Run(c.name, func(t *testing.T) {
			got := c.input.Validate()
			if got != c.exp {
				t.Errorf("got: %v", got)
				t.Errorf("exp: %v", c.exp)
			}
		})
	}
}

func BenchmarkReview_Validate(b *testing.B) {
	for _, c := range reviewCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.input.Validate()
			}
		})
	}
}
