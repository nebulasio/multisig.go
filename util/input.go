package util

import (
    "fmt"
    "strconv"
)

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
