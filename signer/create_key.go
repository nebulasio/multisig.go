package signer

import (
    "encoding/hex"
    "github.com/nebulasio/go-nebulas/crypto/keystore/secp256k1"
    "govote/util"
    "math/rand"
)

const (
    keyFile = "key.txt"
)

func CreateKey() {

    if util.ExistsFile(keyFile) {
        util.PrintError(keyFile, "exists.")
        return
    }

    priKey := make([]byte, 32)
    rand.Read(priKey)
    pubKey, err := secp256k1.GetPublicKey(priKey)
    if err != nil {
        util.PrintError(err)
        return
    }

    hash := util.Sha3256(pubKey)
    content := append([]byte{0x19, 0x57}, util.Rmd160(hash[:])...)
    hash = util.Sha3256(content)
    content = append(content, hash[0:4]...)
    address := util.B58Encode(content)
    if err := util.WriteFile(keyFile, hex.EncodeToString(priKey)+" "+address); err != nil {
        util.PrintError(err)
    } else {
        util.Print("success. ->", keyFile)
    }
}
