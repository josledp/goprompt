package prompt

import (
	"os"
	"testing"
	"time"

	"github.com/josledp/termcolor"

	"github.com/google/go-cmp/cmp"
)

func TestMakeEverteen(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal("error getting current hostname: ", err)
	}
	tt := []struct {
		pr Prompt
		sb string
	}{
		{
			pr: Prompt{
				fullpath: true,
				format:   termcolor.EscapedFormat,
				style:    "Evermeet",
				term: termInfo{
					hostname:   hostname,
					lastrc:     "0",
					pwd:        "~/home",
					user:       "testuser",
					virtualEnv: "venv",
				},
				git: gitInfo{
					branch:        "branch",
					changed:       10,
					untracked:     5,
					stashed:       2,
					staged:        4,
					upstream:      true,
					commitsAhead:  4,
					commitsBehind: 2,
				},

				aws: awsInfo{
					role:   "role:test",
					expire: time.Unix(int64(1506345326), int64(0)),
				},
			},
			sb: "\\[\\033[0m\\]\\[\\033[34m\\](venv) \\[\\033[0m\\]\\[\\033[0m\\]\\[\\033[31m\\]role:test\\[\\033[0m\\]|\\[\\033[0m\\]\\[\\033[1;32m\\]testuser@" + hostname + "\\[\\033[0m\\] \\[\\033[0m\\]\\[\\033[93m\\]0\\[\\033[0m\\] \\[\\033[0m\\]\\[\\033[1;34m\\]~/home\\[\\033[0m\\] \\[\\033[0m\\]\\[\\033[35m\\]branch\\[\\033[0m\\] ↓·2↑·4|\\[\\033[0m\\]\\[\\033[36m\\]●4\\[\\033[0m\\]\\[\\033[0m\\]\\[\\033[36m\\]+10\\[\\033[0m\\]\\[\\033[0m\\]\\[\\033[36m\\]…5\\[\\033[0m\\]\\[\\033[0m\\]\\[\\033[95m\\]⚑2\\[\\033[0m\\]$ ",
		},
		{
			pr: Prompt{
				fullpath: true,
				format:   termcolor.EscapedFormat,
				style:    "Evermeet",
				term: termInfo{
					hostname:   hostname,
					lastrc:     "0",
					pwd:        "~/home",
					user:       "testuser",
					virtualEnv: "venv",
				},
				git: gitInfo{
					branch:        "branch",
					changed:       0,
					untracked:     0,
					stashed:       0,
					staged:        0,
					upstream:      true,
					commitsAhead:  0,
					commitsBehind: 0,
				},

				aws: awsInfo{
					role:   "role:test",
					expire: time.Unix(int64(1506345326), int64(0)),
				},
			},
			sb: "\\[\\033[0m\\]\\[\\033[34m\\](venv) \\[\\033[0m\\]\\[\\033[0m\\]\\[\\033[31m\\]role:test\\[\\033[0m\\]|\\[\\033[0m\\]\\[\\033[1;32m\\]testuser@" + hostname + "\\[\\033[0m\\] \\[\\033[0m\\]\\[\\033[93m\\]0\\[\\033[0m\\] \\[\\033[0m\\]\\[\\033[1;34m\\]~/home\\[\\033[0m\\] \\[\\033[0m\\]\\[\\033[35m\\]branch\\[\\033[0m\\]|\\[\\033[0m\\]\\[\\033[92m\\]✔\\[\\033[0m\\]$ ",
		},
	}
	for _, test := range tt {
		prompt := test.pr.GetPrompt()
		if !cmp.Equal(prompt, test.sb) {
			t.Errorf("error make Everteen prompt %v.\nIt is:     %s\nShould be: %s", test.pr, prompt, test.sb)
		}
	}

}
