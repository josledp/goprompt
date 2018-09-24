package termcolor

import (
	"fmt"
	"strconv"
)

const (
	escapedStart = "\\[\\033["
	escapedEnd   = "m\\]"
	normalStart  = "\033["
	normalEnd    = "m"
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

//GetCode return the term code for modes
func GetCode(modes ...Mode) string {
	return getCode(normalStart, normalEnd, modes...)
}

//GetEscapedCode returns the needed code escaped (mostly for PS1 usage)
func GetEscapedCode(modes ...Mode) string {
	return getCode(escapedStart, escapedEnd, modes...)
}

func getCode(s string, e string, modes ...Mode) string {
	code := s
	for i, m := range modes {
		if i != 0 {
			code += ";"
		}
		code += strconv.Itoa(int(m))
	}
	code += e
	return code
}

//Format returns a strings formated with modes
func Format(s string, modes ...Mode) string {
	return format(s, GetCode(TermReset), GetCode(modes...))
}

//EscapedFormat returns a strings formated with modes escaped with (mostly for PS1 usage)
func EscapedFormat(s string, modes ...Mode) string {
	return format(s, GetEscapedCode(TermReset), GetEscapedCode(modes...))
}

func format(s string, r string, m string) string {
	return fmt.Sprintf("%s%s%s%s", r, m, s, r)
}
