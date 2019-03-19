package tool

import (
    "encoding/json"
    "govote/util"
    "path/filepath"
)

func writeFile(r string, outputPath string) {
    if err := util.WriteFile(outputPath, r); err == nil {
        util.Print("success. ->", outputPath)
    } else {
        util.PrintError(err)
    }
}

func createContractParamItem(data interface{}) interface{} {
    d, err := json.Marshal(data)
    if err != nil {
        util.PrintError(err)
        return nil
    }
    return map[string]interface{}{"data": string(d)}
}

func createSingleContractParam(data interface{}, outputPath string) {
    m := createContractParamItem(data)
    createContractParam([]interface{}{m}, outputPath)
}

func createContractParam(data []interface{}, outputPath string) {
    d, err := json.Marshal(data)
    if err != nil {
        util.PrintError(err)
        return
    }

    a := []interface{}{string(d)}
    r, e := json.Marshal(a)
    if e != nil {
        util.PrintError(e)
        return
    }

    writeFile(string(r), outputPath)
}

func CreateDeleteManagerData(address string) {
    data := map[string]interface{}{"action": "delete_manager", "detail": address}
    output := filepath.Join("output", "delete_manager.json")
    createSingleContractParam(data, output)
}

func CreateAddManagerData(address string) {
    data := map[string]interface{}{"action": "add_manager", "detail": address}
    output := filepath.Join("output", "add_manager.json")
    createSingleContractParam(data, output)
}

func CreateReplaceManagerData(oldAddress string, newAddress string) {
    m := map[string]interface{}{"oldAddress": oldAddress, "newAddress": newAddress}
    data := map[string]interface{}{"action": "replace_manager", "detail": m}
    output := filepath.Join("output", "replace_manager.json")
    createSingleContractParam(data, output)
}

func CreateSendNasData(txsFilePath string) {
    content, err := util.ReadFile(txsFilePath)
    if err != nil {
        util.PrintError(err)
        return
    }

    var txs []interface{}
    err = json.Unmarshal([]byte(content), &txs)
    if err != nil {
        util.PrintError(err)
        return
    }

    var items []interface{}
    for i := 0; i < len(txs); i++ {
        tx := txs[i]
        data := map[string]interface{}{"action": "send_nas", "detail": tx}
        item := createContractParamItem(data)
        if item == nil {
            return
        }
        items = append(items, item)
    }

    outputPath := filepath.Join("output", "send_nas.json")
    createContractParam(items, outputPath)
}

func CreateUpdateSendNasRuleData(ruleFilePath string) {
    content, err := util.ReadFile(ruleFilePath)
    if err != nil {
        util.PrintError(err)
        return
    }

    var rule map[string]interface{}
    err = json.Unmarshal([]byte(content), &rule)
    if err != nil {
        util.PrintError(err)
        return
    }

    m := map[string]interface{}{"action": "update_send_nas_rule", "detail": rule}
    output := filepath.Join("output", "update_send_nas_rule.json")
    createSingleContractParam(m, output)
}

func CreateUpdateSysConfigData(filePath string) {
    content, err := util.ReadFile(filePath)
    if err != nil {
        util.PrintError(err)
        return
    }

    var config map[string]interface{}
    err = json.Unmarshal([]byte(content), &config)
    if err != nil {
        util.PrintError(err)
        return
    }

    m := map[string]interface{}{"action": "update_sys_config", "detail": config}
    output := filepath.Join("output", "update_sys_config.json")
    createSingleContractParam(m, output)
}

func getData(file string) ([]interface{}, error) {
    text, err := util.ReadFile(file)
    if err != nil {
        return nil, err
    }

    var container []interface{}
    err = json.Unmarshal([]byte(text), &container)
    if err != nil {
        return nil, err
    }

    var r []interface{}
    err = json.Unmarshal([]byte(container[0].(string)), &r)
    if err != nil {
        return nil, err
    }
    return r, nil
}

func MergeData(files []string) {
    var r []interface{}
    for _, f := range files {
        data, err := getData(f)
        if err != nil {
            util.PrintError(err)
            return
        }
        r = append(r, data...)
    }
    output := filepath.Join("output", "merge_data.json")
    createContractParam(r, output)
}
