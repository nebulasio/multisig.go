package util

import (
    "errors"
    "io/ioutil"
    "os"
    "path/filepath"
)

func WriteFile(filePath string, content string) error {
    dir := filepath.Dir(filePath)
    if err := CreatDirIfNotExists(dir); err != nil {
        return err
    }
    file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
    if err != nil {
        return err
    }
    defer file.Close()
    if _, err = file.WriteString(content); err != nil {
        return err
    }
    return nil
}

func ReadFile(filePath string) (string, error) {
    if !ExistsFile(filePath) {
        return "", errors.New(filePath + " not found. ")
    }
    bytes, err := ioutil.ReadFile(filePath)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func CreatDirIfNotExists(dir string) error {
    if !ExistsFile(dir) {
        if err := os.MkdirAll(dir, 0777); err != nil {
            return err
        }
    }
    return nil
}

func ExistsFile(file string) bool {
    if _, err := os.Stat(file); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}
