package util

import "strings"

func IsNumber(n string) bool {
    // TODO:
    return true
}

func IsEmptyString(str string) bool {
    return strings.Trim(str, " ") == ""
}
