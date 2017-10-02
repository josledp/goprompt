package plugins

import (
	"os"
	"testing"

	"github.com/josledp/termcolor"
)

func TestPath(t *testing.T) {

	testCases := []struct {
		name           string
		env            []map[string]string
		options        map[string]interface{}
		expectedPwd    string
		expectedPrompt string
	}{
		{
			name: "test_normal_path",
			env: []map[string]string{
				{
					"PWD": "/home",
				},
			},
			expectedPwd:    "/home",
			expectedPrompt: "\\[\\033[0m\\]\\[\\033[1;34m\\]/home\\[\\033[0m\\]",
		},
		{
			name: "test_home",
			env: []map[string]string{
				{
					"PWD":  "/home/test/help",
					"HOME": "/home/test",
				},
			},
			expectedPwd:    "~/help",
			expectedPrompt: "\\[\\033[0m\\]\\[\\033[1;34m\\]~/help\\[\\033[0m\\]",
		},
		{
			name: "test_reduced_path",
			env: []map[string]string{
				{
					"PWD": "/tmp/test",
				},
			},
			options:        map[string]interface{}{"path.fullpath": false},
			expectedPwd:    "test",
			expectedPrompt: "\\[\\033[0m\\]\\[\\033[1;34m\\]test\\[\\033[0m\\]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, env := range tc.env {
				for k, v := range env {
					t.Log(k, v)
					os.Setenv(k, v)
				}
				p := &Path{}
				p.Load(tc.options)
				if p.pwd != tc.expectedPwd {
					t.Fatalf("Pwd do not match:\nGot:      %s\nExpected: %s", p.pwd, tc.expectedPwd)
				}
				pr, _ := p.Get(termcolor.EscapedFormat)
				if pr != tc.expectedPrompt {
					t.Fatalf("Generated prompt do not match:\n%s\n%s", pr, tc.expectedPrompt)
				}
			}
		})
	}

}

func TestPathError(t *testing.T) {

	os.Setenv("PWD", "")

	p := &Path{}
	err := p.Load(map[string]interface{}{})

	if err.Error() != "Unable to get PWD" {
		t.Errorf("Invalid Last command Error: %v", err)
	}

}
