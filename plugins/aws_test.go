package plugins

import (
	"os"
	"testing"
	"time"

	"github.com/josledp/termcolor"
)

func TestAws(t *testing.T) {

	os.Setenv("AWS_ROLE", "test:xx-yy-zz")
	os.Setenv("AWS_SESSION_EXPIRE", "1506345326")
	expected := "\\[\\033[0m\\]\\[\\033[31m\\]test:zz\\[\\033[0m\\]"

	a := &Aws{}
	a.Load(map[string]interface{}{})

	if a.role != "test:zz" {
		t.Error("Invalid AWS role")
	}

	if !(time.Unix(1506345326, int64(0)).Equal(a.expire)) {
		t.Error("Invalid AWS expire time")
	}
	output, _ := a.Get(termcolor.EscapedFormat)
	if output != expected {
		t.Errorf("Expected %s\nGot      %s", expected, output)
	}

}
