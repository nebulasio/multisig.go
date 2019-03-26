package tool

import (
    "encoding/json"
    "path/filepath"
    "strings"
    "vote/util"
)

func verifyAndGetSignees(config map[string]interface{}) string {
    c := util.GetField(config, "signees_count")
    count := int(util.ToNumber(c))

    addresses := util.GetSliceField(config, "addresses")
    if addresses == nil || len(addresses) < count {
        util.PrintError("The number of addresses should be", count)
    }
    for _, a := range addresses {
        util.VerifyAddress(util.ToString(a))
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

func CreateContract(templatePath string, filePath string, output string) {

    text, err := util.ReadFile(filePath)
    if err != nil {
        util.PrintError(err)
    }

    var configs map[string]interface{}
    if err = json.Unmarshal([]byte(text), &configs); err != nil {
        util.PrintError(err)
    }

    config := util.GetMapField(configs, "signees")
    signees := verifyAndGetSignees(config)

    config = util.GetMapField(configs, "constitution")
    constitution := verifyAndGetConstitution(config)

    config = util.GetMapField(configs, "rules")
    rules := verifyAndGetRules(config)

    text, err = util.ReadFile(templatePath)
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
