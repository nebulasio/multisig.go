package main

import (
    "os"
    "strings"
    "vote/signer"
    "vote/tool"
    "vote/util"
)

const helpString = `Command list.

# create key & signature
account create                        Create a private key.
sign <data_file> -o <output_file>     Sign the data file.

# contract
contract gen <contract.conf> -o <output_data_file>                        Create contract.

# data
data-gen <txs_file> -o <output_data_file>                                 Create 'send nas' data.
data-gen remove-signee <address> -o <output_data_file>                    Create 'remove signee' data.
data-gen add-signee <address> -o <output_data_file>                       Create 'add signee' data.
data-gen replace-signee <oldAddress> <newAddress> -o <output_data_file>   Create 'replace signee' data.
data-gen update-rules <send_rules_file> -o <output_data_file>             Create 'update send nas signature rules' data.
data-gen update-constitution <constitution_file> -o <output_data_file>    Create 'update constitution' data.
data-gen merge-file <data_file1> <data_file2> ... -o <output_data_file>   Merge data file.
`

func getArg(index int) string {
    if index > len(os.Args)-1 {
        return ""
    }
    return os.Args[index]
}

func getOutputPath() string {
    for i, a := range os.Args {
        if i > 0 && i < len(os.Args)-1 && strings.ToLower(a) == "-o" {
            return os.Args[i+1]
        }
    }
    return ""
}

func checkArgsLen(length int) {
    if len(os.Args) < length {
        util.PrintError("param error. ")
    }
}

func executeAccountCmd(cmd string) {
    switch cmd {
    case "create":
        signer.CreateKey()
    default:
        util.PrintError("Unknown cmd:", cmd)
    }
}

func executeContractCmd(cmd string) {
    switch cmd {

    case "gen":
        checkArgsLen(4)
        tool.CreateContract(os.Args[3])

    case "deploy":
        // TODO:
        util.PrintError("This feature has not been implemented yet")

    default:
        util.PrintError("Unknown cmd:", cmd)
    }
}

func createData(dataType string) {
    output := getOutputPath()

    switch dataType {

    case util.ActionRemoveSignee:
        checkArgsLen(4)
        tool.CreateDeleteManagerData(os.Args[3], output)

    case util.ActionAddSignee:
        checkArgsLen(4)
        tool.CreateAddManagerData(os.Args[3], output)

    case util.ActionReplaceSignee:
        checkArgsLen(5)
        tool.CreateReplaceManagerData(os.Args[3], os.Args[4], output)

    case util.ActionUpdateRules:
        checkArgsLen(4)
        tool.CreateUpdateSendNasRuleData(os.Args[3], output)

    case util.ActionUpdateConstitution:
        checkArgsLen(4)
        tool.CreateUpdateSysConfigData(os.Args[3], output)

    case "merge-file":
        checkArgsLen(5)
        var files []string
        for index, item := range os.Args {
            if index > 2 {
                if strings.ToLower(item) == "-o" {
                    break
                }
                files = append(files, item)
            }
        }
        if len(files) < 2 {
            util.PrintError("Need at least two files. ")
        }
        tool.MergeData(files, output)

    default:
        checkArgsLen(3)
        tool.CreateSendNasData(os.Args[2], output)
    }
}

func executeCmd(cmd string) {
    switch strings.ToLower(cmd) {

    case "account":
        checkArgsLen(3)
        executeAccountCmd(strings.ToLower(os.Args[2]))

    case "contract":
        checkArgsLen(3)
        executeContractCmd(strings.ToLower(os.Args[2]))

    case "data-gen":
        checkArgsLen(3)
        createData(strings.ToLower(os.Args[2]))

    case "sign":
        checkArgsLen(3)
        keyPath := getArg(3)
        if strings.ToLower(keyPath) == "-o" {
            keyPath = ""
        }
        signer.Sign(os.Args[2], keyPath, getOutputPath())

    default:
        util.Print("Unknown cmd:", cmd)
    }
}

func main() {
    if len(os.Args) == 1 {
        util.Print(helpString)
        return
    }
    executeCmd(os.Args[1])
}
