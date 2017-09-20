package colors

import (
	"fmt"
	"strconv"
)

const (
	start = "\033["
	end   = "m"
)

//Mode type is for defining term modes
type Mode int

// Constans borrowed from github.com/fatih/color
// Base attributes
const (
	TermReset Mode = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack Mode = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Mode = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Mode = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack Mode = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

//Reset returns the term code for reset colors
func Reset() string {
	return GetCode(TermReset)
}

//GetCode return the term code for modes
func GetCode(modes ...Mode) string {
	code := start
	for i, m := range modes {
		if i != 0 {
			code += ";"
		}
		code += strconv.Itoa(int(m))
	}
	code += end
	return code
}

//Format returns a strings formated with modes
func Format(s string, modes ...Mode) string {
	return fmt.Sprintf("%s%s%s%s", Reset(), GetCode(modes...), s, Reset())
}
