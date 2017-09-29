package plugins

import (
	"os"
	"strconv"
	"strings"
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

//Load is the load function of the plugin
func (a *Aws) Load(options map[string]interface{}) error {
	role := os.Getenv("AWS_ROLE")
	if role != "" {
		tmp := strings.Split(role, ":")
		role = tmp[0]
		tmp = strings.Split(tmp[1], "-")
		role += ":" + tmp[2]
	}
	a.role = role
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
		} else if d < 1800 {
			t = termcolor.FgBlue
		} else if d < 600 {
			t = termcolor.FgYellow
		}
		return format(a.role, t), []termcolor.Mode{t}
	}
	return "", nil
}
