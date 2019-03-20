package signer

import (
    "encoding/hex"
    "errors"
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
    sig, err = util.Sign([]byte(data.(string)), bytesKey)
    if err != nil {
        return err
    }

    strSig := string(sig)
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
    for _, data := range array {
        util.VerifyData(data)
        if err := sign(data, key); err != nil {
            util.PrintError(err)
        }
    }
    util.SerializeDataListToFile(array, outputPath)
}
