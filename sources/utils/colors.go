package utils

import (
	"github.com/fatih/color"
)

func Faint(message string) string {
	return color.New(color.Faint).SprintFunc()(message)
}

func Success(message string) string {
	return color.New(color.FgGreen).SprintFunc()(message)
}

func Failure(message string) string {
	return color.New(color.FgRed).SprintFunc()(message)
}

func ClearLine() string {
	return "\033[A\033[K"
}
