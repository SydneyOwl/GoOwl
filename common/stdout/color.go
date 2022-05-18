package stdout

import (
	"github.com/sydneyowl/GoOwl/common/global"

	"github.com/fatih/color"
)

//Return string in red
func Red(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgRed).SprintFunc()(msg)
	}
	return msg
}

//Return string in green
func Green(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgGreen).SprintFunc()(msg)
	}
	return msg
}

//Return string in yellow
func Yellow(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgYellow).SprintFunc()(msg)
	}
	return msg
}

//Return string in blue
func Blue(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgBlue).SprintFunc()(msg)
	}
	return msg
}

//Return string in Magenta
func Magenta(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgHiMagenta).SprintFunc()(msg)
	}
	return msg
}

//Return string in cyan
func Cyan(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgCyan).SprintFunc()(msg)
	}
	return msg
}

//Return string in white
func White(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgWhite).SprintFunc()(msg)
	}
	return msg
}
