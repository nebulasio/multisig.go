package util

func verifyVoteAction(action map[string]interface{}) {
    name := GetStringField(action, "name")
    detail := GetMapField(action, "detail")

    switch name {
    case "callContract":
        verifyCallContractAction(detail)
    default:
        PrintError("Unknown action name:", name)
    }
}

func verifyCallContractAction(detail map[string]interface{}) {
    address := GetStringField(detail, "address")
    _ = GetStringField(detail, "func")
    args := VerifyAndGetField(detail, "args")
    if !IsSlice(args) {
        PrintError("args is not array.")
    }
    VerifyAddress(address)
}
