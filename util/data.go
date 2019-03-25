package util

import (
    "encoding/json"
    "fmt"
)

func writeDataToFile(r string, outputPath string) {
    if err := WriteFile(outputPath, r); err == nil {
        Print("success. ->", outputPath)
    } else {
        PrintError(err)
    }
}

func CreateContractData(data interface{}) map[string]interface{} {
    d, err := json.Marshal(data)
    if err != nil {
        PrintError(err)
    }
    return map[string]interface{}{"data": string(d)}
}

func SerializeDataToFile(data interface{}, outputPath string) {
    m := CreateContractData(data)
    SerializeDataListToFile([]map[string]interface{}{m}, outputPath)
}

func SerializeDataListToFile(data []map[string]interface{}, outputPath string) {
    r, err := json.Marshal(data)
    if err != nil {
        PrintError(err)
    }
    writeDataToFile(string(r), outputPath)
}

func DeserializeDataFile(file string) ([]map[string]interface{}, error) {
    text, err := ReadFile(file)
    if err != nil {
        return nil, err
    }

    var r []map[string]interface{}
    err = json.Unmarshal([]byte(text), &r)
    if err != nil {
        return nil, err
    }
    return r, nil
}

func DeserializeData(data string) map[string]interface{} {
    var r map[string]interface{}
    err := json.Unmarshal([]byte(data), &r)
    if err != nil {
        PrintError("Data error. ")
    }
    return r
}

func VerifyData(data map[string]interface{}) (string, string) {
    action := GetStringField(data, "action")
    detail := VerifyAndGetField(data, "detail")

    v := ""
    switch action {
    case ActionRemoveSignee:
        VerifyAddress(ToString(detail))

    case ActionAddSignee:
        VerifyAddress(ToString(detail))

    case ActionReplaceSignee:
        verifyReplaceManagerData(ToMap(detail))

    case ActionUpdateRules:
        VerifySendRules(ToMap(detail))

    case ActionUpdateConstitution:
        VerifyConstitution(ToMap(detail))

    case ActionSend:
        v = verifySendNasData(ToMap(detail))

    case ActionVote:
        v = verifyVoteData(ToMap(detail))

    default:
        PrintError("Action", action, "is not supported.")
    }
    return action, v
}

func verifyReplaceManagerData(data map[string]interface{}) {
    oldAddress := GetStringField(data, "oldAddress")
    newAddress := GetStringField(data, "newAddress")

    VerifyAddress(oldAddress)
    VerifyAddress(newAddress)
    if oldAddress == newAddress {
        PrintError("Data error. ")
    }
}

func verifySendNasData(item map[string]interface{}) string {
    id := GetNotEmptyStringField(item, "id")
    to := GetStringField(item, "to")
    VerifyAddress(to)
    value := GetStringField(item, "value")
    VerifyNumber(value)
    return id
}

func verifyVoteData(item map[string]interface{}) string {
    id := GetNotEmptyStringField(item, "id")
    _ = GetStringField(item, "content")

    p := GetStringField(item, "proportionOfApproved")
    VerifyProportions(p)

    action := VerifyAndGetField(item, "approvedAction")
    verifyVoteAction(ToMap(action))

    return id
}

func VerifyConstitution(data map[string]interface{}) {
    ver := GetStringField(data, "version")
    VerifyNumber(ver)

    p := GetMapField(data, "proportionOfSigners")
    ks := []interface{}{"updateConstitution", "updateSendRules", "addSignee", "removeSignee", "replaceSignee", "vote"}
    n := 0
    for k, v := range p {
        if Contains(ks, k) {
            n++
        }
        VerifyProportions(ToString(v))
    }
    if n != len(ks) {
        PrintError("Constitution data error. ")
    }
}

func VerifySendRules(data map[string]interface{}) {
    ver := GetStringField(data, "version")
    VerifyNumber(ver)

    rules := GetSliceField(data, "rules")
    if len(rules) <= 0 {
        PrintError("rules is empty. ")
    }

    v := 0.0
    for _, i := range rules {
        r := ToMap(i)

        p := GetStringField(r, "proportionOfSigners")
        VerifyProportions(p)

        t := GetStringField(r, "startValue")
        startValue := ParseFloat(t)

        if v == -1 || startValue != v {
            PrintError("Rules error. ", startValue, v)
        }

        e := GetStringField(r, "endValue")
        if e != Infinity {
            v = ParseFloat(e)
            if startValue >= v {
                PrintError("Rules error. ")
            }
        } else {
            v = -1
        }
    }
    if v != -1 {
        PrintError("Rules error. ")
    }
}

func PrintData(data map[string]interface{}) {
    action, _ := data["action"]
    fmt.Println("\n============", action, "============")
    detail, _ := data["detail"]
    d, _ := json.MarshalIndent(detail, "", "  ")
    fmt.Println(string(d))
}
