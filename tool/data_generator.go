package tool

import (
    "encoding/json"
    "path/filepath"
    "vote/util"
)

func CreateDeleteManagerData(address string, output string) {
    data := map[string]interface{}{"action": util.ActionRemoveSignee, "detail": address}
    util.VerifyData(data)
    if util.IsEmptyString(output) {
        output = filepath.Join("output", "remove-signee.json")
    }
    util.SerializeDataToFile(data, output)
}

func CreateAddManagerData(address string, output string) {
    data := map[string]interface{}{"action": util.ActionAddSignee, "detail": address}
    util.VerifyData(data)
    if util.IsEmptyString(output) {
        output = filepath.Join("output", "add-signee.json")
    }
    util.SerializeDataToFile(data, output)
}

func CreateReplaceManagerData(oldAddress string, newAddress string, output string) {
    m := map[string]interface{}{"oldAddress": oldAddress, "newAddress": newAddress}
    data := map[string]interface{}{"action": util.ActionReplaceSignee, "detail": m}
    util.VerifyData(data)
    if util.IsEmptyString(output) {
        output = filepath.Join("output", "replace-signee.json")
    }
    util.SerializeDataToFile(data, output)
}

func CreateSendNasData(txsFilePath string, output string) {
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
    ids := make([]interface{}, 10)
    for i := 0; i < len(txs); i++ {
        tx := txs[i]
        data := map[string]interface{}{"action": util.ActionSend, "detail": tx}
        _, id := util.VerifyData(data)
        if util.Contains(ids, id) {
            util.PrintError("tx.id", id, "has been repeated. ")
        }
        ids = append(ids, id)
        item := util.CreateContractData(data)
        if item == nil {
            return
        }
        items = append(items, item)
    }

    if util.IsEmptyString(output) {
        output = filepath.Join("output", "send.json")
    }
    util.SerializeDataListToFile(items, output)
}

func CreateVoteData(voteFilePath string, output string) {
    content, err := util.ReadFile(voteFilePath)
    if err != nil {
        util.PrintError(err)
        return
    }

    var votes []interface{}
    err = json.Unmarshal([]byte(content), &votes)
    if err != nil {
        util.PrintError(err)
        return
    }

    var items []map[string]interface{}
    ids := make([]interface{}, 10)
    for i := 0; i < len(votes); i++ {
        vote := votes[i]
        data := map[string]interface{}{"action": util.ActionVote, "detail": vote}
        _, id := util.VerifyData(data)
        if util.Contains(ids, id) {
            util.PrintError("vote.id", id, "has been repeated. ")
        }
        ids = append(ids, id)
        item := util.CreateContractData(data)
        if item == nil {
            return
        }
        items = append(items, item)
    }

    if util.IsEmptyString(output) {
        output = filepath.Join("output", "vote.json")
    }
    util.SerializeDataListToFile(items, output)
}

func CreateUpdateSendNasRuleData(ruleFilePath string, output string) {
    content, err := util.ReadFile(ruleFilePath)
    if err != nil {
        util.PrintError(err)
    }

    var rules map[string]interface{}
    err = json.Unmarshal([]byte(content), &rules)
    if err != nil {
        util.PrintError(err)
    }

    data := map[string]interface{}{"action": util.ActionUpdateRules, "detail": rules}
    util.VerifyData(data)
    if util.IsEmptyString(output) {
        output = filepath.Join("output", "update-rules.json")
    }
    util.SerializeDataToFile(data, output)
}

func CreateUpdateSysConfigData(filePath string, output string) {
    content, err := util.ReadFile(filePath)
    if err != nil {
        util.PrintError(err)
    }

    var config map[string]interface{}
    err = json.Unmarshal([]byte(content), &config)
    if err != nil {
        util.PrintError(err)
    }

    data := map[string]interface{}{"action": util.ActionUpdateConstitution, "detail": config}
    util.VerifyData(data)
    if util.IsEmptyString(output) {
        output = filepath.Join("output", "update-constitution.json")
    }
    util.SerializeDataToFile(data, output)
}

func MergeData(files []string, output string) {
    var r []map[string]interface{}
    for _, f := range files {
        array, err := util.DeserializeDataFile(f)
        if err != nil {
            util.PrintError(err)
        }
        for _, container := range array {
            str, ok := container["data"]
            if !ok {
                util.PrintError("Data error. ")
            } else {
                data := util.DeserializeData(str.(string))
                util.VerifyData(data)
            }
        }
        r = append(r, array...)
    }
    if util.IsEmptyString(output) {
        output = filepath.Join("output", "merge-file.json")
    }
    util.SerializeDataListToFile(r, output)
}
