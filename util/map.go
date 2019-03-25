package util

func VerifyAndGetField(dict map[string]interface{}, fieldName string) interface{} {
    r, ok := dict[fieldName]
    if !ok {
        PrintError(fieldName, " field is empty. ")
    }
    return r
}

func GetStringField(dict map[string]interface{}, fieldName string) string {
    r := VerifyAndGetField(dict, fieldName)
    return ToString(r)
}

func GetNotEmptyStringField(dict map[string]interface{}, fieldName string) string {
    s := GetStringField(dict, fieldName)
    if IsEmptyString(s) {
        PrintError(fieldName, " field is empty. ")
    }
    return s
}

func GetSliceField(dict map[string]interface{}, fieldName string) []interface{} {
    r := VerifyAndGetField(dict, fieldName)
    return ToSlice(r)
}

func GetMapField(dict map[string]interface{}, fieldName string) map[string]interface{} {
    r := VerifyAndGetField(dict, fieldName)
    return ToMap(r)
}
