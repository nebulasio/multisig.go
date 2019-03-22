package tool

import (
    "encoding/json"
    "path/filepath"
    "strings"
    "vote/util"
)

func verifyAndGetSignees(config map[string]interface{}) string {
    c, ok := config["signees_count"]
    if !ok {
        util.PrintError("No fields found for 'signees_count'. ")
    }
    count := int(c.(float64))

    array, ok := config["addresses"]
    if !ok {
        util.PrintError("No fields found for 'addresses'. ")
    }

    addresses := array.([]interface{})
    if addresses == nil || len(addresses) < count {
        util.PrintError("The number of addresses should be", count)
    }

    for _, a := range addresses {
        util.VerifyAddress(a.(string))
    }

    if data, err := json.Marshal(addresses); err == nil {
        return string(data)
    } else {
        util.PrintError(err)
        return ""
    }
}

func verifyAndGetConstitution(config map[string]interface{}) string {
    util.VerifyConstitution(config)
    if data, err := json.Marshal(config); err == nil {
        return string(data)
    } else {
        util.PrintError(err)
        return ""
    }
}

func verifyAndGetRules(config map[string]interface{}) string {
    util.VerifySendRules(config)
    if data, err := json.Marshal(config); err == nil {
        return string(data)
    } else {
        util.PrintError(err)
        return ""
    }
}

func CreateContract(filePath string, output string) {

    text, err := util.ReadFile(filePath)
    if err != nil {
        util.PrintError(err)
    }

    var configs map[string]map[string]interface{}
    if err = json.Unmarshal([]byte(text), &configs); err != nil {
        util.PrintError(err)
    }

    config, ok := configs["signees"]
    if !ok {
        util.PrintError("'signees' field cannot be empty")
    }
    signees := verifyAndGetSignees(config)

    config, ok = configs["constitution"]
    if !ok {
        util.PrintError("'constitution' field cannot be empty")
    }
    constitution := verifyAndGetConstitution(config)

    config, ok = configs["rules"]
    if !ok {
        util.PrintError("'rules' field cannot be empty")
    }
    rules := verifyAndGetRules(config)

    text, err = util.ReadFile("template.js")
    if err != nil {
        util.PrintError(err)
    }
    text = strings.Replace(text, "SIGNEES", signees, -1)
    text = strings.Replace(text, "CONSTITUTION", constitution, -1)
    text = strings.Replace(text, "RULES", rules, -1)
    if util.IsEmptyString(output) {
        output = filepath.Join("output", "contract.js")
    }
    if err = util.WriteFile(output, text); err != nil {
        util.PrintError(err)
    }
    util.Print("success. ->", output)
}
