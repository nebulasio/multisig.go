package tool

import (
    "encoding/json"
    "govote/util"
    "path/filepath"
    "strings"
)

func CreateContract(filePath string) {

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

    managerCount, contains := info["manager_count"]
    if !contains {
        util.PrintError("No fields found for 'manager_count'. ")
        return
    }
    count := int(managerCount.(float64))

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

    outputPath := filepath.Join("output", "contract.js")
    if data, err := json.Marshal(addrs); err == nil {
        text = strings.Replace(text, "MANAGERS", string(data), -1)
        if err := util.WriteFile(outputPath, text); err != nil {
            util.PrintError(err)
            return
        }
    }

    util.Print("success. ->", outputPath)
}
