package color

import (
	"fmt"
)

// using zap color
type Color func(string) string
type color uint8

const (
	black color = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

func (c color) draw(s string) string {
	// return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
	return fmt.Sprintf("\033[%dm%s\033[0m", uint8(c), s)
}

// Color represents a text color.

// Add adds the coloring to the given string.

// Yellow ...
func Yellow(msg string) string {
	return yellow.draw(msg)
}

// Red ...
func Red(msg string) string {
	return red.draw(msg)
}

// Blue ...
func Blue(msg string) string {
	return blue.draw(msg)
}

// Green ...
func Green(msg string) string {
	return green.draw(msg)
}
