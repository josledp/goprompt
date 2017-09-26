package prompt

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
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal("error getting current hostname: ", err)
	}
	tt := []struct {
		environ map[string]string
		ti      termInfo
	}{
		{
			map[string]string{"PWD": "/", "HOSTNAME": hostname, "HOME": "/home/test1", "USER": "test1", "LAST_COMMAND_RC": "0", "VIRTUAL_ENV": ""},
			termInfo{pwd: "/", hostname: hostname, user: "test1", lastrc: "0", virtualEnv: ""},
		},
		{
			map[string]string{"PWD": "/home/test1", "HOSTNAME": hostname, "HOME": "/home/test1", "USER": "test1", "LAST_COMMAND_RC": "127", "VIRTUAL_ENV": "/virtualenv/testvenv"},
			termInfo{pwd: "~", hostname: hostname, user: "test1", lastrc: "127", virtualEnv: "testvenv"},
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
