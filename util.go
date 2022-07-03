package main

import (
	"regexp"
)

func stripANSI(s string) string {
	ansi := `[\x1b\x9b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]`
	re := regexp.MustCompile(ansi)
	return re.ReplaceAllString(s, "")
}
