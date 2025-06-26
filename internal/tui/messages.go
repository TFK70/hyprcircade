package tui

import (
	"fmt"
)

func GetCompletedMsg(msg string) string {
	return fmt.Sprintf("%s %s", " ", msg)
}
