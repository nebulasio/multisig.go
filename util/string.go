package util

import (
    "github.com/btcsuite/btcutil/base58"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

const (
    ActionRemoveSignee       = "remove-signee"
    ActionAddSignee          = "add-signee"
    ActionReplaceSignee      = "replace-signee"
    ActionUpdateRules        = "update-rules"
    ActionUpdateConstitution = "update-constitution"
    ActionSend               = "send"
    ActionVote               = "vote"

    Infinity = "INFINITY"
)

var (
    VotingValues = []string{"agree", "disagree", "abstain"}
)

func IsEmptyString(str string) bool {
    return strings.Trim(str, " ") == ""
}

func VerifyNumber(str string) {
    reg := regexp.MustCompile(`^\d+(\.\d+)?$`)
    if !reg.Match([]byte(str)) {
        PrintError("Data error. ")
    }
}

func VerifyAddress(address string) {
    if len(address) != 35 || strings.Index(address, "n1") != 0 {
        PrintError(address, "is not a valid nas address.")
    }
    content := base58.Decode(address)
    l := len(content)
    hash := Sha3256(content[:l-4])
    if !reflect.DeepEqual(hash[:4], content[l-4:]) {
        PrintError(address, "is not a valid nas address.")
    }
}

func VerifyProportions(str string) {
    VerifyNumber(str)
    p := ParseFloat(str)
    if p <= 0 || p > 1 {
        PrintError("Proportion error")
    }
}

func ParseFloat(str string) float64 {
    f, err := strconv.ParseFloat(str, 64)
    if err != nil {
        PrintError(err)
    }
    return f
}
