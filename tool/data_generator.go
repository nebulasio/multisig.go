package tool

import (
    "encoding/json"
    "path/filepath"
    "vote/util"
)

func CreateDeleteManagerData(address string) {
    data := map[string]interface{}{"action": util.ActionDeleteManager, "detail": address}
    util.VerifyData(data)
    output := filepath.Join("output", "delete_manager.json")
    util.SerializeDataToFile(data, output)
}

func CreateAddManagerData(address string) {
    data := map[string]interface{}{"action": util.ActionAddManager, "detail": address}
    util.VerifyData(data)
    output := filepath.Join("output", "add_manager.json")
    util.SerializeDataToFile(data, output)
}

func CreateReplaceManagerData(oldAddress string, newAddress string) {
    m := map[string]interface{}{"oldAddress": oldAddress, "newAddress": newAddress}
    data := map[string]interface{}{"action": util.ActionReplaceManager, "detail": m}
    util.VerifyData(data)
    output := filepath.Join("output", "replace_manager.json")
    util.SerializeDataToFile(data, output)
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

    var items []map[string]interface{}
    for i := 0; i < len(txs); i++ {
        tx := txs[i]
        data := map[string]interface{}{"action": util.ActionSendNas, "detail": tx}
        util.VerifyData(data)
        item := util.CreateContractData(data)
        if item == nil {
            return
        }
        items = append(items, item)
    }

    outputPath := filepath.Join("output", "send_nas.json")
    util.SerializeDataListToFile(items, outputPath)
}

func CreateUpdateSendNasRuleData(ruleFilePath string) {
    content, err := util.ReadFile(ruleFilePath)
    if err != nil {
        util.PrintError(err)
    }

    var rule map[string]interface{}
    err = json.Unmarshal([]byte(content), &rule)
    if err != nil {
        util.PrintError(err)
    }

    data := map[string]interface{}{"action": util.ActionUpdateSendNasRule, "detail": rule}
    util.VerifyData(data)
    output := filepath.Join("output", "update_send_nas_rule.json")
    util.SerializeDataToFile(data, output)
}

func CreateUpdateSysConfigData(filePath string) {
    content, err := util.ReadFile(filePath)
    if err != nil {
        util.PrintError(err)
    }

    var config map[string]interface{}
    err = json.Unmarshal([]byte(content), &config)
    if err != nil {
        util.PrintError(err)
    }

    data := map[string]interface{}{"action": util.ActionUpdateSysConfig, "detail": config}
    util.VerifyData(data)
    output := filepath.Join("output", "update_sys_config.json")
    util.SerializeDataToFile(data, output)
}

func MergeData(files []string) {
    var r []map[string]interface{}
    for _, f := range files {
        array, err := util.DeserializeDataFile(f)
        if err != nil {
            util.PrintError(err)
        }
        for _, data := range array {
            util.VerifyData(data)
        }
        r = append(r, array...)
    }
    output := filepath.Join("output", "merge_data.json")
    util.SerializeDataListToFile(r, output)
}
