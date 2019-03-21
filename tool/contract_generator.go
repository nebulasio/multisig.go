package tool

import (
    "encoding/json"
    "path/filepath"
    "strings"
    "vote/util"
)

func CreateContract(filePath string, output string) {

    text, err := util.ReadFile(filePath)
    if err != nil {
        util.PrintError(err)
        return
    }

    var info map[string]interface{}
    if err = json.Unmarshal([]byte(text), &info); err != nil {
        util.PrintError(err)
        return
    }

    signeesCount, contains := info["signees_count"]
    if !contains {
        util.PrintError("No fields found for 'signees_count'. ")
        return
    }
    count := int(signeesCount.(float64))

    addresses, containsAddress := info["addresses"]
    if !containsAddress {
        util.PrintError("No fields found for 'addresses'. ")
        return
    }

    addrs := addresses.([]interface{})
    if addrs == nil || len(addrs) < count {
        util.PrintError("The number of addresses should be", count)
        return
    }

    text, err = util.ReadFile("template.js")
    if err != nil {
        util.PrintError(err)
        return
    }

    if util.IsEmptyString(output) {
        output = filepath.Join("output", "contract.js")
    }
    if data, err := json.Marshal(addrs); err == nil {
        text = strings.Replace(text, "SIGNEES", string(data), -1)
        if err := util.WriteFile(output, text); err != nil {
            util.PrintError(err)
            return
        }
    }

    util.Print("success. ->", output)
}
