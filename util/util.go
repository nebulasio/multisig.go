package util

import (
    "fmt"
    "github.com/btcsuite/btcutil/base58"
    "github.com/nebulasio/go-nebulas/crypto/keystore/secp256k1"
    "golang.org/x/crypto/ripemd160"
    "golang.org/x/crypto/sha3"
    "os"
)

func Sha3256(data []byte) [32]byte {
    return sha3.Sum256(data)
}

func Rmd160(data []byte) []byte {
    h := ripemd160.New()
    h.Write(data)
    return h.Sum(nil)
}

func B58Encode(data []byte) string {
    return base58.Encode(data)
}

func Sign(data []byte, secKey []byte) ([]byte, error) {
    return secp256k1.Sign(data, secKey)
}

func RunActions(actions []func()) {
    for _, f := range actions {
        f()
    }
}

func Print(msg ...interface{}) {
    fmt.Println(msg...)
    fmt.Println()
}

func PrintError(msg ...interface{}) {
    fmt.Print("[Error] ")
    fmt.Println(msg...)
    fmt.Println()
    os.Exit(0)
}

func IsNumber(o interface{}) bool {
    switch o.(type) {
    case float64:
        return true
    }
    return false
}

func IsString(o interface{}) bool {
    switch o.(type) {
    case string:
        return true
    }
    return false
}

func IsSlice(o interface{}) bool {
    switch o.(type) {
    case []interface{}:
        return true
    }
    return false
}

func IsMap(o interface{}) bool {
    switch o.(type) {
    case map[string]interface{}:
        return true
    }
    return false
}

func ToNumber(o interface{}) float64 {
    if !IsNumber(o) {
        PrintError(o, "is not number. ")
    }
    return o.(float64)
}

func ToString(o interface{}) string {
    if !IsString(o) {
        PrintError(o, "is not map. ")
    }
    return o.(string)
}

func ToSlice(o interface{}) []interface{} {
    if !IsSlice(o) {
        PrintError(o, "is not array. ")
    }
    return o.([]interface{})
}

func ToMap(o interface{}) map[string]interface{} {
    if !IsMap(o) {
        PrintError(o, "is not map. ")
    }
    return o.(map[string]interface{})
}
