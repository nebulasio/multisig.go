package util

func GetField(dict map[string]interface{}, fieldName string) interface{} {
    r, ok := dict[fieldName]
    if !ok {
        PrintError(fieldName, " field is empty. ")
    }
    return r
}

func GetStringField(dict map[string]interface{}, fieldName string) string {
    return ToString(GetField(dict, fieldName))
}

func GetNotEmptyStringField(dict map[string]interface{}, fieldName string) string {
    s := GetStringField(dict, fieldName)
    if IsEmptyString(s) {
        PrintError(fieldName, " field is empty. ")
    }
    return s
}

func GetSliceField(dict map[string]interface{}, fieldName string) []interface{} {
    return ToSlice(GetField(dict, fieldName))
}

func GetMapField(dict map[string]interface{}, fieldName string) map[string]interface{} {
    return ToMap(GetField(dict, fieldName))
}
