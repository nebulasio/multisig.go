package util

import "encoding/json"

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
    d, err := json.Marshal(data)
    if err != nil {
        PrintError(err)
    }

    a := []string{string(d)}
    r, e := json.Marshal(a)
    if e != nil {
        PrintError(e)
    }

    writeDataToFile(string(r), outputPath)
}

func DeserializeDataFile(file string) ([]map[string]interface{}, error) {
    text, err := ReadFile(file)
    if err != nil {
        return nil, err
    }

    var container []string
    err = json.Unmarshal([]byte(text), &container)
    if err != nil {
        return nil, err
    }

    var r []map[string]interface{}
    err = json.Unmarshal([]byte(container[0]), &r)
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

func VerifyData(data map[string]interface{}) {
    action, ok := data["action"]
    if !ok {
        PrintError("Data error.")
    }

    detail, ok := data["detail"]
    if !ok {
        PrintError("Data error.")
    }

    switch action {
    case ActionRemoveSignee:
        VerifyAddress(detail.(string))

    case ActionAddSignee:
        VerifyAddress(detail.(string))

    case ActionReplaceSignee:
        verifyReplaceManagerData(detail.(map[string]interface{}))

    case ActionSend:
        verifySendNasData(detail.([]map[string]interface{}))

    case ActionUpdateRules:
        verifySendNasRule(detail.(map[string]interface{}))

    case ActionUpdateConstitution:
        verifySysConfig(detail.(map[string]interface{}))

    default:
        PrintError("Action", action, "is not supported.")
    }
}

func verifyReplaceManagerData(data map[string]interface{}) {
    oldAddress, ok := data["oldAddress"]
    if !ok {
        PrintError("oldAddress is empty. ")
    }

    newAddress, ok := data["newAddress"]
    if !ok {
        PrintError("newAddress is empty. ")
    }

    VerifyAddress(oldAddress.(string))
    VerifyAddress(newAddress.(string))
}

func verifySendNasData(data []map[string]interface{}) {
    ids := make([]interface{}, 10)
    for _, item := range data {
        id, ok := item["id"]
        if !ok || IsEmptyString(id.(string)) {
            PrintError("tx.id is empty. ")
        }
        if Contains(ids, id) {
            PrintError("tx.id", id, "has been repeated. ")
        }
        ids = append(ids, id.(string))
        to, ok := item["to"]
        if !ok {
            PrintError("tx.to is empty. ")
        }
        VerifyAddress(to.(string))
        value, ok := item["value"]
        if !ok {
            PrintError("tx.value is empty. ")
        }
        VerifyNumber(value.(string))
    }
}

func verifySysConfig(data map[string]interface{}) {
    ver, ok := data["version"]
    if !ok {
        PrintError("version is empty. ")
    }
    VerifyNumber(ver.(string))

    t, ok := data["proportionOfSigners"]
    if !ok {
        PrintError("proportionOfSigners is empty. ")
    }

    p := t.(map[string]string)
    ks := []interface{}{"updateSysConfig", "updateSendNasRule", "addManager", "deleteManager", "replaceManager"}
    n := 0
    for k, v := range p {
        if Contains(ks, k) {
            n++
        }
        VerifyProportions(v)
    }
    if n != len(ks) {
        PrintError("sys config data error. ")
    }
}

func verifySendNasRule(data map[string]interface{}) {
    ver, ok := data["version"]
    if !ok {
        PrintError("version is empty. ")
    }
    VerifyNumber(ver.(string))

    t, ok := data["rules"]
    if !ok {
        PrintError("rules is empty. ")
    }

    rules := t.([]map[string]interface{})
    if len(rules) <= 0 {
        PrintError("rules is empty. ")
    }

    v := 0.0
    for _, r := range rules {
        p, ok := r["proportionOfSigners"]
        if !ok {
            PrintError("proportionOfSigners is empty. ")
        }
        VerifyProportions(p.(string))

        t, ok := r["startValue"]
        if !ok {
            PrintError("startValue is empty. ")
        }
        startValue := ParseFloat(t.(string))
        if v == -1 || startValue != v {
            PrintError("Rules error. ")
        }

        e, ok := r["endValue"]
        if !ok {
            PrintError("endValue is empty. ")
        }
        var endValue float64
        if e != Infinity {
            endValue := ParseFloat(e.(string))
            if startValue >= endValue {
                PrintError("Rules error. ")
            }
        } else {
            endValue = -1
        }
        v = endValue
    }
    if v != -1 {
        PrintError("Rules error. ")
    }
}
