package main

import (
	"regexp"
	"sort"
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

// returns slice of sorted keys from map[string]T
func sortKeys(m map[string]App) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// insert element into slice at given index
func insertAtIndex(index int, element string, array []string) []string {
	array = append(array[:index+1], array[index:]...)
	array[index] = element
	return array
}

// strip all ANSI escape characters
func stripANSI(s string) string {
	ansi := `[\x1b\x9b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]`
	re := regexp.MustCompile(ansi)
	return re.ReplaceAllString(s, "")
}
