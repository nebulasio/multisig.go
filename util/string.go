package util

import (
    "regexp"
    "strconv"
    "strings"
)

const (
    ActionRemoveSignee       = "remove-signee"
    ActionAddSignee          = "add-signee"
    ActionReplaceSignee      = "replace-signee"
    ActionSend               = "send"
    ActionUpdateRules        = "update-rules"
    ActionUpdateConstitution = "update-constitution"

    Infinity = "INFINITY"
)

func VerifyNumber(str string) {
    reg := regexp.MustCompile(`^\d+(\.\d+)?$`)
    if !reg.Match([]byte(str)) {
        PrintError("Data error. ")
    }
}

func ParseFloat(str string) float64 {
    f, err := strconv.ParseFloat(str, 64)
    if err != nil {
        PrintError(err)
    }
    return f
}

func VerifyProportions(str string) {
    VerifyNumber(str)
    p := ParseFloat(str)
    if p <= 0 || p > 1 {
        PrintError("Proportion error")
    }
}

func IsEmptyString(str string) bool {
    return strings.Trim(str, " ") == ""
}
