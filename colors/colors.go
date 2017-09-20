package colors

import "fmt"

const (
	start = "\033["
	end   = "m"
)

var modesMap map[string]string

//Reset returns the term code for reset colors
func Reset() string {
	return GetCode("reset")
}

//GetCode return the term code for modes
func GetCode(modes ...string) string {
	code := start
	for i, m := range modes {
		if i != 0 {
			code += ";"
		}
		code += modesMap[m]
	}
	code += end
	return code
}

//Format returns a strings formated with modes
func Format(s string, modes ...string) string {
	return fmt.Sprintf("%s%s%s%s", Reset(), GetCode(modes...), s, Reset())
}

func init() {
	modesMap = make(map[string]string)
	modesMap["reset"] = "0"
	modesMap["bold"] = "1"
	modesMap["underline"] = "4"
	modesMap["boldOff"] = "21"
	modesMap["underlineOff"] = "24"
	modesMap["black"] = "30"
	modesMap["red"] = "31"
	modesMap["green"] = "32"
	modesMap["yellow"] = "33"
	modesMap["blue"] = "34"
	modesMap["magenta"] = "35"
	modesMap["cyan"] = "36"
	modesMap["white"] = "37"
}
