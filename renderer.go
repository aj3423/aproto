package aproto

import (
	"fmt"
	"html"

	"github.com/fatih/color"
)

type Renderer interface {
	INDENT() string       // indent string
	NEWLINE() string      // renders a new line
	IDTYPE(string) string // full id+type hex bytes
	ID(string) string     // only id number
	TYPE(string) string   // only type string
	NUM(string) string    // a number
	STR(string) string    // a string
}

// --------------- Console ---------------
type ConsoleRenderer struct {
}

func (r *ConsoleRenderer) INDENT() string {
	return "    "
}
func (r *ConsoleRenderer) NEWLINE() string {
	return "\n"
}
func (r *ConsoleRenderer) IDTYPE(s string) string {
	return color.HiRedString(s)
}
func (r *ConsoleRenderer) ID(s string) string {
	return color.HiGreenString(s)
}
func (r *ConsoleRenderer) TYPE(s string) string {
	return color.YellowString(s)
}
func (r *ConsoleRenderer) NUM(s string) string {
	return color.HiCyanString(s)
}
func (r *ConsoleRenderer) STR(s string) string {
	return color.HiYellowString(s)
}

// --------------- HTML ---------------
func wrap_tag(s, tag string) string {
	return fmt.Sprintf("<%s>%s</%s>", tag, s, tag)
}

type HtmlRenderer struct {
}

func (r *HtmlRenderer) INDENT() string {
	return "&nbsp;&nbsp;&nbsp;&nbsp;"
}
func (r *HtmlRenderer) NEWLINE() string {
	return "</br>"
}
func (r *HtmlRenderer) IDTYPE(s string) string {
	return fmt.Sprintf("<font color='#ff2200'>%s</font>", s)
}
func (r *HtmlRenderer) ID(s string) string {
	return fmt.Sprintf("<font color='#00ff11'>%s</font>", s)
}
func (r *HtmlRenderer) TYPE(s string) string {
	return fmt.Sprintf("<font color='#808000'>%s</font>", s)
}
func (r *HtmlRenderer) NUM(s string) string {
	return fmt.Sprintf("<font color='cyan'>%s</font>", s)
}
func (r *HtmlRenderer) STR(s string) string {
	return fmt.Sprintf("<font color='yellow'>%s</font>", html.EscapeString(s))
}
