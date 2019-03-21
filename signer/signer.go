package signer

import (
    "encoding/hex"
    "errors"
    "fmt"
    "path/filepath"
    "strings"
    "vote/util"
)

func readKey(keyPath string) (string, error) {
    text, err := util.ReadFile(keyPath)
    if err != nil {
        return "", err
    }
    key := strings.Split(text, " ")[0]
    if len(key) != 64 {
        return "", errors.New("key error. ")
    }
    return key, nil
}

func sign(info map[string]interface{}, key string) error {
    data, ok := info["data"]
    if !ok {
        return errors.New("data error. ")
    }

    bytesKey, err := hex.DecodeString(key)
    if err != nil {
        return err
    }

    var sig []byte
    hash := util.Sha3256([]byte(data.(string)))
    sig, err = util.Sign(hash[:], bytesKey)
    if err != nil {
        return err
    }

    strSig := hex.EncodeToString(sig)
    _, c := info["sigs"]
    var ss []interface{}
    if !c {
        ss = []interface{}{}
        info["sigs"] = ss
    } else {
        ss = info["sigs"].([]interface{})
    }
    if util.Contains(ss, strSig) {
        return errors.New("you have signed. ")
    }
    ss = append(ss, strSig)
    info["sigs"] = ss
    return nil
}

func vote(info map[string]interface{}, key string, voteValue string) error {
    data, ok := info["data"]
    if !ok {
        return errors.New("data error. ")
    }

    bytesKey, err := hex.DecodeString(key)
    if err != nil {
        return err
    }

    var sig []byte
    hash := util.Sha3256([]byte(data.(string) + voteValue))
    sig, err = util.Sign(hash[:], bytesKey)
    if err != nil {
        return err
    }

    strSig := hex.EncodeToString(sig)
    _, ok = info["votes"]
    var ss []interface{}
    if !ok {
        ss = []interface{}{}
        info["votes"] = ss
    } else {
        ss = info["votes"].([]interface{})
    }

    for _, v := range ss {
        d := v.(map[string]string)
        s, ok := d["sig"]
        if !ok {
            return errors.New("Sig error. ")
        }
        if s == strSig {
            return errors.New("you have signed. ")
        }
    }
    info["votes"] = append(ss, map[string]string{"sig": strSig, "value": voteValue})
    return nil
}

func Sign(filePath string, keyPath string, outputPath string) {
    if util.IsEmptyString(filePath) {
        util.PrintError("data file path is empty. ")
    }

    if util.IsEmptyString(keyPath) {
        keyPath = "key.txt"
    }

    if util.IsEmptyString(outputPath) {
        outputPath = filepath.Join("output", filepath.Base(filePath))
    }

    key, err := readKey(keyPath)
    if err != nil {
        util.PrintError(err)
    }

    array, err := util.DeserializeDataFile(filePath)
    if err != nil {
        util.PrintError(err)
    }
    sendIds := make([]interface{}, 0, 10)
    voteIds := make([]interface{}, 0, 10)
    sigActions := make([]func(), 0, 3)
    votesActions := make([]func(), 0, 3)
    printSigDatas := make([]func(), 0, 10)
    for _, item := range array {
        d, ok := item["data"]
        if !ok {
            util.PrintError("Data error.")
        }
        data := util.DeserializeData(d.(string))
        action, r := util.VerifyData(data)
        if action == util.ActionSend {
            if util.Contains(sendIds, r) {
                util.PrintError("tx.id", r, "has been repeated. ")
            }
            sendIds = append(sendIds, r)
        } else if action == util.ActionVote {
            if util.Contains(voteIds, r) {
                util.PrintError("vote.id", r, "has been repeated. ")
            }
            voteIds = append(voteIds, r)
        }

        tempItem := item
        if action == util.ActionVote {
            votesActions = append(votesActions, func() {
                util.PrintData(data)
                fmt.Println()
                voteValue := util.GetVoteResult()
                if err := vote(tempItem, key, voteValue); err != nil {
                    util.PrintError(err)
                }
                fmt.Println()
            })
        } else {
            printSigDatas = append(printSigDatas, func() {
                util.PrintData(data)
            })
            sigActions = append(sigActions, func() {
                if err := sign(tempItem, key); err != nil {
                    util.PrintError(err)
                }
            })
        }
    }
    util.RunActions(votesActions)
    util.RunActions(printSigDatas)
    fmt.Println()
    if util.AgreeSig() {
        util.RunActions(sigActions)
        util.Print()
    } else {
        util.Print()
        util.PrintError("You canceled signing the data.")
    }

    util.SerializeDataListToFile(array, outputPath)
}
