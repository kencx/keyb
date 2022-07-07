package util

import (
	"regexp"
)

func StripANSI(s string) string {
	ansi := `[\x1b\x9b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]`
	re := regexp.MustCompile(ansi)
	return re.ReplaceAllString(s, "")
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
