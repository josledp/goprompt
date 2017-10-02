package plugins

import (
	"os"
	"testing"

	"github.com/josledp/termcolor"
)

func TestUserChar(t *testing.T) {

	testCases := []struct {
		user     string
		expected string
	}{
		{
			user:     "testuser",
			expected: "$",
		},
		{
			user:     "root",
			expected: "#",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.user, func(t *testing.T) {
			os.Setenv("USER", tc.user)
			uc := &UserChar{}
			uc.Load(nil)

			if uc.user != tc.user {
				t.Error("Invalid user")
			}

			output, _ := uc.Get(termcolor.EscapedFormat)
			if output != tc.expected {
				t.Errorf("Expected %s\nGot      %s", tc.expected, output)
			}
		})
	}

}
