// monkey/utils/utils.go

package utils

import (
    "fmt"
    "bytes"
    "reflect"
)

func IsNill(tmp any) bool {
    if tmp == nil || reflect.ValueOf(tmp).IsNil() {
        return true
    }
    return false
}

func HasInput(x any) (string, bool) {
    var value = reflect.ValueOf(x)
    if value.Kind() != reflect.Struct {  // is struct ?
        return "", false
    }
    var f = value.FieldByName("input")
    if !f.IsValid() {
        return "", false
    }
    return f.String(), true
}

func GetInput[T any](args []T) (string, error) {
    var out bytes.Buffer
    for _, x := range args {
        var val, isValid = HasInput(x)
        if !isValid {
            return "", fmt.Errorf("Found one struct without the 'input' field on it")
        }
        out.WriteString(val + ";\n")
    }
    return out.String(), nil
}
