package util

import "regexp"

var numericRegex = regexp.MustCompile(`[^0-9 ]+`)

func ClearString(str string) string {
	return numericRegex.ReplaceAllString(str, "")
}
