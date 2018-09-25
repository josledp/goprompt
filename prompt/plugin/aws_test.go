package plugin

import (
	"os"
	"testing"
	"time"

	"github.com/josledp/termcolor"
)

func TestAws(t *testing.T) {

	os.Setenv("AWS_ROLE", "test:xx-yy-zz")
	os.Setenv("AWS_SESSION_EXPIRE", "1506345326")
	expectedPrompt := "\\[\\033[0m\\]\\[\\033[31m\\]test:xx-yy-zz\\[\\033[0m\\]"
	expectedRole := "test:xx-yy-zz"

	a := &Aws{}
	a.Load(nil)

	if a.role != expectedRole {
		t.Error("Expected role %s, got %s", expectedRole, a.role)
	}

	if !(time.Unix(1506345326, int64(0)).Equal(a.expire)) {
		t.Error("AWS expire time error. expected %d, got %d", 1506345326, a.expire)
	}
	output, _ := a.Get(termcolor.EscapedFormat)
	if output != expectedPrompt {
		t.Errorf("Expected %s\nGot      %s", expectedPrompt, output)
	}

}
