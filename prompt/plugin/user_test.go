package plugin

import (
	"os"
	"testing"

	"github.com/josledp/termcolor"
)

func TestUser(t *testing.T) {

	testCases := []struct {
		user     string
		expected string
	}{
		{
			user:     "testuser",
			expected: "\\[\\033[0m\\]\\[\\033[1;32m\\]testuser\\[\\033[0m\\]",
		},
		{
			user:     "root",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.user, func(t *testing.T) {
			os.Setenv("USER", tc.user)
			u := &User{}
			u.Load(nil)

			if u.user != tc.user {
				t.Error("Invalid user")
			}

			output, _ := u.Get(termcolor.EscapedFormat)
			if output != tc.expected {
				t.Errorf("Expected %s\nGot      %s", tc.expected, output)
			}
		})
	}

}
