package util

import (
    "fmt"
    "github.com/btcsuite/btcutil/base58"
    "github.com/nebulasio/go-nebulas/crypto/keystore/secp256k1"
    "golang.org/x/crypto/ripemd160"
    "golang.org/x/crypto/sha3"
    "os"
    "reflect"
    "strconv"
    "strings"
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

func VerifyAddress(address string) {
    if len(address) != 35 || strings.Index(address, "n1") != 0 {
        PrintError(address, "is not a valid nas address.")
    }
    content := base58.Decode(address)
    l := len(content)
    hash := Sha3256(content[:l-4])
    if !reflect.DeepEqual(hash[:4], content[l-4:]) {
        PrintError(address, "is not a valid nas address.")
    }
}

func GetInput(desc string) string {
    fmt.Print(desc)
    var r string
    _, _ = fmt.Scanf("%s", &r)
    return r
}

func AgreeSig() bool {
    r := GetInput("Confirm the signature of the above data? \n1: ok  2: cancel\nEnter value: ")
    for !Contains([]interface{}{"1", "2"}, r) {
        Print("\nPlease enter 1 or 2. ")
        r = GetInput("1: ok  2: cancel \nEnter value: ")
    }
    return r == "1"
}

func GetVoteResult() string {
    r := GetInput("Voting values. \n1: agree  2: disagree  3: abstain\nEnter value: ")
    for !Contains([]interface{}{"1", "2", "3"}, r) {
        Print("\nPlease enter 1 or 2 or 3. ")
        r = GetInput("1: agree  2: disagree  3: abstain\nEnter value: ")
    }
    i, _ := strconv.Atoi(r)
    return VotingValues[i-1]
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
