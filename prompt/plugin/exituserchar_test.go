package plugin

import (
	"os"
	"testing"

	"github.com/josledp/termcolor"
)

func TestExitUserChar(t *testing.T) {

	testCases := []struct {
		name           string
		lastrc         string
		user           string
		expectedPrompt string
	}{
		{
			name:           "normaluser_noerror",
			lastrc:         "0",
			user:           "normaluser",
			expectedPrompt: "$",
		},
		{
			name:           "normaluser_error",
			lastrc:         "10",
			user:           "normaluser",
			expectedPrompt: "\\[\\033[0m\\]\\[\\033[91m\\]$\\[\\033[0m\\]",
		},
		{
			name:           "root_noerror",
			lastrc:         "0",
			user:           "root",
			expectedPrompt: "#",
		},
		{
			name:           "root_error",
			lastrc:         "1",
			user:           "root",
			expectedPrompt: "\\[\\033[0m\\]\\[\\033[91m\\]#\\[\\033[0m\\]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("LAST_COMMAND_RC", tc.lastrc)
			os.Setenv("USER", tc.user)
			euc := &ExitUserChar{}
			euc.Load(nil)
			pr, _ := euc.Get(termcolor.EscapedFormat)
			if pr != tc.expectedPrompt {
				t.Fatalf("Generated prompt do not match:\n%s\n%s", pr, tc.expectedPrompt)
			}

		})
	}

}
