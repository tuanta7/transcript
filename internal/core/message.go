package core

import (
	"fmt"
	"strings"
)

func ErrorMessage(msg string) string {
	return fmt.Sprintf("Error: %s", msg)
}

func WarningMessage(msg string) string {
	return fmt.Sprintf("Warning: %s", msg)
}

func IsErrorMessage(msg string) bool {
	return strings.HasPrefix(msg, ErrorMessage(""))
}

func IsWarningMessage(msg string) bool {
	return strings.HasPrefix(msg, WarningMessage(""))
}
