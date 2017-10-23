package plugins

import (
	"os"
	"testing"

	"github.com/josledp/termcolor"
)

func TestHostname(t *testing.T) {
	host, err := os.Hostname()
	if err != nil {
		t.Fatal("Unable to get hostname")
	}
	testCases := []struct {
		user     string
		expected string
	}{
		{
			user:     "testuser",
			expected: "\\[\\033[0m\\]\\[\\033[1;32m\\]" + host + "\\[\\033[0m\\]",
		},
		{
			user:     "root",
			expected: "\\[\\033[0m\\]\\[\\033[1;31m\\]" + host + "\\[\\033[0m\\]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.user, func(t *testing.T) {
			os.Setenv("USER", tc.user)
			h := &Hostname{}
			h.Load(nil)

			if h.user != tc.user {
				t.Error("Invalid user")
			}

			if h.hostname != host {
				t.Error("Invalid host")
			}

			output, _ := h.Get(termcolor.EscapedFormat)
			if output != tc.expected {
				t.Errorf("Expected %s\nGot      %s", tc.expected, output)
			}
		})
	}

}
