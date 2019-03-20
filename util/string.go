package util

import (
    "regexp"
    "strconv"
    "strings"
)

const (
    ActionDeleteManager     = "delete_manager"
    ActionAddManager        = "add_manager"
    ActionReplaceManager    = "replace_manager"
    ActionSendNas           = "send_nas"
    ActionUpdateSendNasRule = "update_send_nas_rule"
    ActionUpdateSysConfig   = "update_sys_config"

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
