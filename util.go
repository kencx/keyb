package main

import (
	"regexp"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func insertAtIndex(index int, element string, array []string) []string {
	array = append(array[:index+1], array[index:]...)
	array[index] = element
	return array
}

func stripANSI(s string) string {
	ansi := `[\x1b\x9b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]`
	re := regexp.MustCompile(ansi)
	return re.ReplaceAllString(s, "")
}
