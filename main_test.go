package main

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAwsInfo(t *testing.T) {
	os.Setenv("AWS_ROLE", "test:xx-yy-zz")
	os.Setenv("AWS_SESSION_EXPIRE", "1506345326")

	ai := getAwsInfo()

	if ai.role != "test:zz" {
		t.Error("Invalid AWS role")
	}

	if !(time.Unix(1506345326, int64(0)).Equal(ai.expire)) {
		t.Error("Invalid AWS expire time")
	}
}

func TestTermInfo(t *testing.T) {
	tt := []struct {
		environ map[string]string
		ti      termInfo
	}{
		{
			map[string]string{"PWD": "/", "HOSTNAME": "hosttest", "HOME": "/home/test1", "USER": "test1", "LAST_COMMAND_RC": "0", "VIRTUAL_ENV": ""},
			termInfo{pwd: "/", hostname: "hosttest", user: "test1", lastrc: "0", virtualEnv: ""},
		},
		{
			map[string]string{"PWD": "/home/test1", "HOSTNAME": "hosttest", "HOME": "/home/test1", "USER": "test1", "LAST_COMMAND_RC": "127", "VIRTUAL_ENV": "/virtualenv/testvenv"},
			termInfo{pwd: "~", hostname: "hosttest", user: "test1", lastrc: "127", virtualEnv: "testvenv"},
		},
	}

	for _, test := range tt {
		for k, v := range test.environ {
			os.Setenv(k, v)
		}
		ti := getTermInfo()
		ti.virtualEnv = getPythonVirtualEnv()

		if !cmp.Equal(test.ti, ti, cmp.AllowUnexported(termInfo{}, awsInfo{}, gitInfo{})) {
			t.Error("Invalid case detected: ", test, ti)
		}
	}
}

func TestMakeEverteen(t *testing.T) {
	pi := promptInfo{
		term: termInfo{
			hostname:   "hosttest",
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
	}
	shouldbe := "\\[\\033[0m\\]\\[\\033[34m\\](venv) \\[\\033[0m\\][\\[\\033[0m\\]\\[\\033[31m\\]role:test\\[\\033[0m\\]|\\[\\033[0m\\]\\[\\033[1;32m\\]hosttest\\[\\033[0m\\] \\[\\033[0m\\]\\[\\033[93m\\]0\\[\\033[0m\\] \\[\\033[0m\\]\\[\\033[1;34m\\]~/home\\[\\033[0m\\] \\[\\033[0m\\]\\[\\033[35m\\]branch\\[\\033[0m\\] ↓·2↑·4|\\[\\033[0m\\]\\[\\033[36m\\]●4\\[\\033[0m\\]\\[\\033[0m\\]\\[\\033[36m\\]+10\\[\\033[0m\\]\\[\\033[0m\\]\\[\\033[36m\\]…5\\[\\033[0m\\]\\[\\033[0m\\]\\[\\033[95m\\]⚑2\\[\\033[0m\\]]$"
	prompt := makePromptEvermeet(pi)
	if !cmp.Equal(prompt, shouldbe) {
		t.Error(prompt)
		t.Error(shouldbe)
	}
}
