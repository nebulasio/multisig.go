package util

import (
    "fmt"
    "github.com/btcsuite/btcutil/base58"
    "github.com/nebulasio/go-nebulas/crypto/keystore/secp256k1"
    "golang.org/x/crypto/ripemd160"
    "golang.org/x/crypto/sha3"
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

func VerifyAddress(address string) error {
    // TODO:
    return nil
}

func Print(msg ...interface{}) {
    fmt.Println(msg...)
    fmt.Println()
}

func PrintError(msg ...interface{}) {
    fmt.Print("[Error] ")
    fmt.Println(msg...)
    fmt.Println()
}
