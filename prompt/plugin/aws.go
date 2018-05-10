package plugin

import (
	"os"
	"strconv"
	"time"

	"github.com/josledp/termcolor"
)

//Aws is the plugin struct
type Aws struct {
	role   string
	expire time.Time
}

//Name returns the plugin name
func (Aws) Name() string {
	return "aws"
}

//Help returns help information about this plugin
func (Aws) Help() (description string, options map[string]string) {
	description = "This plugins show aws information(it needs AWS_ROLE + AWS_SESSION_EXPIRE non standard environment variables"
	return
}

//Load is the load function of the plugin
func (a *Aws) Load(Prompter) error {
	a.role = os.Getenv("AWS_ROLE")
	iExpire, _ := strconv.ParseInt(os.Getenv("AWS_SESSION_EXPIRE"), 10, 0)
	a.expire = time.Unix(iExpire, int64(0))
	return nil
}

//Get returns the string to use in the prompt
func (a Aws) Get(format func(string, ...termcolor.Mode) string) (string, []termcolor.Mode) {
	if a.role != "" {
		t := termcolor.FgGreen
		d := time.Until(a.expire).Seconds()
		if d < 0 {
			t = termcolor.FgRed
		} else if d < 600 {
			t = termcolor.FgYellow
		} else if d < 1800 {
			t = termcolor.FgBlue
		}
		return format(a.role, t), nil
	}
	return "", nil
}
