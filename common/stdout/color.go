package stdout

import (
	"github.com/sydneyowl/GoOwl/common/global"

	"github.com/fatih/color"
)

// Red return string in red
func Red(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgRed).SprintFunc()(msg)
	}
	return msg
}

// Green return string in green
func Green(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgGreen).SprintFunc()(msg)
	}
	return msg
}

// Yellow return string in yellow
func Yellow(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgYellow).SprintFunc()(msg)
	}
	return msg
}

// Blue return string in blue
func Blue(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgBlue).SprintFunc()(msg)
	}
	return msg
}

// Magenta return string in Magenta
func Magenta(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgHiMagenta).SprintFunc()(msg)
	}
	return msg
}

// Cyan return string in cyan
func Cyan(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgCyan).SprintFunc()(msg)
	}
	return msg
}

// White return string in white
func White(msg string) string {
	if global.OS == "linux" {
		return color.New(color.FgWhite).SprintFunc()(msg)
	}
	return msg
}
