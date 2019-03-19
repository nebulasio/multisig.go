package util

func Contains(array []interface{}, item interface{}) bool {
    for i := 0; i < len(array); i++ {
        if array[i] == item {
            return true
        }
    }
    return false
}
