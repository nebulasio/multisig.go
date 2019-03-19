package main

import (
    "govote/signer"
    "govote/tool"
    "govote/util"
    "os"
    "strings"
)

const helpString = `Command list.

# create key & signature
create             Create a private key.
sign <file_path>   Sign the data file.

# contract
contract <contract_config_file_path>                 Create contract.

# data
data delete <address>                                Create 'delete manager' data.
data add <address>                                   Create 'add manager' data.
data replace <oldAddress> <newAddress>               Create 'replace manager' data.
data send_nas <txs_file_path>                        Create 'send nas' data.
data update_send_nas_rule <send_nas_rule_file_path>  Create 'update send nas signature rules' data.
data update_sys_config <sys_config_file_path>        Create 'update sys config' data.
data merge <data_file1_path> <data_file2_path> ...   Merge data file.
`

func checkArgsLen(length int) bool {
    if len(os.Args) < length {
        util.PrintError("param error. ")
        return false
    }
    return true
}

func createData(dataType string) {
    switch dataType {
    case "delete":
        if checkArgsLen(4) {
            tool.CreateDeleteManagerData(os.Args[3])
        }
    case "add":
        if checkArgsLen(4) {
            tool.CreateAddManagerData(os.Args[3])
        }
    case "replace":
        if checkArgsLen(5) {
            tool.CreateReplaceManagerData(os.Args[3], os.Args[4])
        }
    case "send_nas":
        if checkArgsLen(4) {
            tool.CreateSendNasData(os.Args[3])
        }
    case "update_send_nas_rule":
        if checkArgsLen(4) {
            tool.CreateUpdateSendNasRuleData(os.Args[3])
        }
    case "update_sys_config":
        if checkArgsLen(4) {
            tool.CreateUpdateSysConfigData(os.Args[3])
        }
    case "merge":
        if checkArgsLen(5) {
            var files []string
            for index, item := range os.Args {
                if index > 2 {
                    files = append(files, item)
                }
            }
            tool.MergeData(files)
        }
    }
}

func executeCmd(cmd string) {
    switch strings.ToLower(strings.Trim(cmd, " ")) {
    case "create":
        signer.CreateKey()
    case "sign":
        if checkArgsLen(3) {
            signer.Sign(os.Args[2], "", "")
        }
    case "contract":
        if checkArgsLen(3) {
            tool.CreateContract(os.Args[2])
        }
    case "data":
        if checkArgsLen(3) {
            createData(os.Args[2])
        }
    default:
        util.Print("unknown cmd:", cmd)
    }
}

func main() {
    if len(os.Args) == 1 {
        util.Print(helpString)
        return
    }
    executeCmd(os.Args[1])
}
