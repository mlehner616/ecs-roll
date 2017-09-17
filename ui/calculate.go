package ui

import (
	"strings"
)

func getHeight(s string) int {
	if len(s) == 0 {
		return 0
	}

	return strings.Count(s, "\n") + 1
}

func getWidth(s string) int {
	lines := strings.Split(s, "\n")

	max := 0
	for _, line := range lines {
		len := len(line)

		if len > max {
			max = len
		}
	}

	return max
}
