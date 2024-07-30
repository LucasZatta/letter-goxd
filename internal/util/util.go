package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var numericRegex = regexp.MustCompile(`[^0-9]+`)

func ClearString(str string) string {
	return strings.Trim(numericRegex.ReplaceAllString(str, ""), " ")
}

func StringElementToInt(str string) int {
	num, err := strconv.Atoi(ClearString(str))
	if err != nil {
		fmt.Println(err)
	}
	return num
}
